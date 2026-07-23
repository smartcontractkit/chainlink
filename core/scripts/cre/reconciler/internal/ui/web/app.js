// CRE Reconciler — Web UI
(function () {
    'use strict';

    const state = {
        nodes: [],
        chains: [],
        bootstrap: '',
        gateways: [],
        gatewayNodes: [],
        namespace: '',
        chartDir: '',
        dons: [],
        capabilities: [],
        capabilityConfigs: {},
        selectedDON: null,
        infra: { type: 'griddle', chartValues: '', namespace: '' },
        jd: { grpc: 'grpc-job-distributor.main.stage.cldev.sh:443', domain: 'cre', environment: 'dev', useTLS: true },
    };

    // --- API helpers ---
    async function api(url, opts) {
        const resp = await fetch(url, opts);
        if (!resp.ok) {
            const body = await resp.json().catch(() => ({}));
            throw new Error(body.error || resp.statusText);
        }
        return resp;
    }

    async function apiGet(url) {
        const resp = await api(url);
        return resp.json();
    }

    async function apiPost(url, data) {
        const resp = await api(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });
        return resp.json().catch(() => ({}));
    }

    // --- Toast ---
    function toast(msg, type) {
        const el = document.getElementById('toast');
        const msgEl = document.getElementById('toastMsg');
        msgEl.textContent = msg;
        msgEl.className = 'px-4 py-3 rounded-lg shadow-lg text-sm ' +
            (type === 'error' ? 'bg-red-600 text-white' :
                type === 'success' ? 'bg-green-600 text-white' :
                    'bg-gray-800 text-white');
        el.classList.remove('hidden');
        setTimeout(() => el.classList.add('hidden'), 3000);
    }

    // --- Tab switching ---
    function switchTab(tabName) {
        document.querySelectorAll('.tab-content').forEach(el => el.classList.add('hidden'));
        document.querySelectorAll('[data-tab]').forEach(el => {
            el.className = el.dataset.tab === tabName ? 'tab-active' : 'tab-inactive';
        });
        document.getElementById('tab-' + tabName).classList.remove('hidden');

        if (tabName === 'status') loadState();
        if (tabName === 'configs') renderCapConfigs();
        if (tabName === 'chains') renderChains();
    }

    // --- Load nodes from chart values ---
    async function loadNodes() {
        try {
            const data = await apiGet('/api/nodes');
            state.nodes = data.nodes || [];
            state.bootstrap = data.bootstrap || '';
            state.gateways = data.gateways || [];
            state.namespace = data.namespace || '';
            state.chartDir = data.chartDir || '';
            document.getElementById('namespaceLabel').textContent =
                state.namespace ? `Namespace: ${state.namespace}` : '';

            // Prepopulate config fields if not already set by desired state
            if (!state.infra.namespace) state.infra.namespace = state.namespace;
            if (!state.infra.chartValues) state.infra.chartValues = state.chartDir;

            state.gateways.forEach(gwName => {
                if (!state.gatewayNodes.find(gwn => gwn.node === gwName)) {
                    state.gatewayNodes.push({ node: gwName, don: '' });
                }
            });
        } catch (e) {
            toast('Failed to load nodes: ' + e.message, 'error');
        }
    }

    // --- Load desired state ---
    async function loadDesired() {
        try {
            const data = await apiGet('/api/desired');
            state.chains = (data.chains || []).map(c => ({
                chainId: c.chainId,
                wsUrl: c.wsUrl || '',
                httpUrl: c.httpUrl || '',
                registry: c.registry || false,
            }));
            state.dons = (data.dons || []).map(d => ({
                name: d.name,
                donTypes: d.donTypes || [],
                capabilities: d.capabilities || [],
                nodes: d.nodes || [],
                bootstrapNode: d.bootstrapNode || '',
                exposesRemoteCapabilities: d.exposesRemoteCapabilities || false,
                registryBasedLaunchAllowlist: d.registryBasedLaunchAllowlist || [],
                capabilityConfigs: d.capabilityConfigs || {},
            }));
            state.infra = data.infra || { type: 'griddle', chartValues: '', namespace: state.namespace };
            state.jd = data.jd || {};
            // Pre-fill defaults if not set
            if (!state.jd.grpc) state.jd.grpc = 'grpc-job-distributor.main.stage.cldev.sh:443';
            if (!state.jd.domain) state.jd.domain = 'cre';
            if (!state.jd.environment) state.jd.environment = 'dev';
            if (state.jd.useTLS === undefined) state.jd.useTLS = true;
            state.capabilityConfigs = data.capabilityConfigs || {};
            state.gatewayNodes = (data.gatewayNodes || []).map(gwn => ({ node: gwn.node, don: gwn.don }));

            // Ensure infra.namespace is set
            if (!state.infra.namespace) state.infra.namespace = state.namespace;

            renderDONs();
            renderConfig();
            renderChains();
        } catch (e) {
            // No desired state yet — start fresh
            renderDONs();
            renderConfig();
            renderChains();
        }
    }

    // --- Chains ---

    // chainVariantKey identifies a distinct (chainId, wsUrl, httpUrl) tuple —
    // the same chain ID can legitimately show up more than once with different
    // RPC URLs (e.g. a gateway node pointed at a different, sometimes unusable,
    // RPC than worker nodes), and every variant must survive so the user can
    // prune the wrong one instead of one silently overwriting the other.
    function chainVariantKey(c) {
        return `${c.chainId}|${c.wsUrl}|${c.httpUrl}`;
    }

    // loadDiscoveredChains preloads the Chains tab from the chart's node config.
    // Called automatically on startup, and again from the "Refresh from chart"
    // button (silent suppresses the toast for the automatic call).
    async function loadDiscoveredChains(silent) {
        try {
            const data = await apiGet('/api/chains/discovered');
            const discovered = data.chains || [];
            const existingKeys = new Set(state.chains.map(chainVariantKey));
            let added = 0;
            discovered.forEach(c => {
                const entry = { chainId: c.chainId, wsUrl: c.wsUrl || '', httpUrl: c.httpUrl || '', registry: false };
                const key = chainVariantKey(entry);
                if (existingKeys.has(key)) return;
                state.chains.push(entry);
                existingKeys.add(key);
                added++;
            });
            renderChains();
            if (!silent) {
                toast(added > 0 ? `Loaded ${added} chain(s) from chart` : 'No new chains found in chart', added > 0 ? 'success' : 'info');
            }
        } catch (e) {
            toast('Failed to load chains from chart: ' + e.message, 'error');
        }
    }

    function renderChains() {
        const container = document.getElementById('chainsContent');
        if (!container) return;

        if (state.chains.length === 0) {
            container.innerHTML = '<div class="text-xs text-gray-400 text-center py-4">No chains found in the chart — add [[EVM]] network config to a node in the chart, then click "Refresh from chart"</div>';
            return;
        }

        container.innerHTML = state.chains.map((c, idx) => `
        <div class="grid grid-cols-12 gap-2 items-center mb-2" data-chain-row="${idx}">
            <input type="number" value="${c.chainId || ''}" placeholder="chain_id"
                class="col-span-2 text-xs px-2 py-1 border border-gray-200 rounded font-mono"
                onchange="updateChainField(${idx}, 'chainId', this.value)">
            <input type="text" value="${c.wsUrl}" placeholder="ws_url"
                class="col-span-4 text-xs px-2 py-1 border border-gray-200 rounded font-mono"
                onchange="updateChainField(${idx}, 'wsUrl', this.value)">
            <input type="text" value="${c.httpUrl}" placeholder="http_url"
                class="col-span-4 text-xs px-2 py-1 border border-gray-200 rounded font-mono"
                onchange="updateChainField(${idx}, 'httpUrl', this.value)">
            <label class="col-span-1 flex items-center justify-center text-[10px] text-gray-500">
                <input type="radio" name="registryChain" ${c.registry ? 'checked' : ''}
                    onchange="setRegistryChain(${idx})" class="mr-1">
                registry
            </label>
            <button onclick="removeChain(${idx})" class="col-span-1 text-gray-300 hover:text-red-500 text-sm">&times;</button>
        </div>
    `).join('');
    }

    window.updateChainField = function (idx, field, value) {
        if (field === 'chainId') {
            state.chains[idx][field] = parseInt(value) || 0;
        } else {
            state.chains[idx][field] = value;
        }
    };

    window.setRegistryChain = function (idx) {
        state.chains.forEach((c, i) => { c.registry = i === idx; });
        renderChains();
    };

    window.removeChain = function (idx) {
        state.chains.splice(idx, 1);
        renderChains();
    };

    // --- Load capabilities catalog ---
    async function loadCapabilities() {
        try {
            const data = await apiGet('/api/capabilities');
            state.capabilities = data.capabilities || [];
            // Prepopulate capability configs with defaults if not already set
            if (data.defaults && Object.keys(state.capabilityConfigs).length === 0) {
                state.capabilityConfigs = data.defaults;
            }
            renderCapabilityCatalog();
        } catch (e) {
            toast('Failed to load capabilities: ' + e.message, 'error');
        }
    }

    // --- Render: DON columns ---
    // DON membership is read-only in the UI — it is always derived server-side
    // from the chart's don-name label (see /api/desired), never edited here.
    function renderDONs() {
        const container = document.getElementById('donColumns');
        container.innerHTML = state.dons.map((don, idx) => donColumnHTML(don, idx)).join('');

        renderCapabilityCatalog();
        renderGatewayAssignments();
    }

    function renderGatewayAssignments() {
        const container = document.getElementById('gatewayAssignments');
        if (!container) return;

        const gatewayDONs = state.dons.filter(d => d.donTypes.includes('gateway'));
        const workflowDONs = state.dons.filter(d => d.donTypes.includes('workflow'));

        if (gatewayDONs.length === 0) {
            container.innerHTML = '';
            container.classList.add('hidden');
            return;
        }

        container.classList.remove('hidden');
        let html = `
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
            <h2 class="text-sm font-semibold text-gray-700">Gateway DON Assignments</h2>
            <p class="text-xs text-gray-400 mb-3">Map each gateway DON to the workflow DON it serves</p>
    `;

        gatewayDONs.forEach(don => {
            const assignedNode = don.nodes.find(n => state.gateways.includes(n)) || don.nodes[0] || '';
            const assignment = state.gatewayNodes.find(gwn => gwn.node === assignedNode);
            const selectedDON = assignment ? assignment.don : '';

            html += `
            <div class="gateway-assignment">
                <span class="text-sm font-medium text-gray-800">${don.name}</span>
                <select class="don-type-select" onchange="updateGatewayDON('${don.name}', this.value)">
                    <option value="">— Select workflow DON —</option>
                    ${workflowDONs.map(wf => `
                        <option value="${wf.name}" ${selectedDON === wf.name ? 'selected' : ''}>${wf.name}</option>
                    `).join('')}
                </select>
            </div>
        `;
        });

        html += '</div>';
        container.innerHTML = html;
    }

    window.updateGatewayDON = function (gatewayDONName, workflowDONName) {
        const don = state.dons.find(d => d.name === gatewayDONName);
        if (!don) return;
        don.nodes.forEach(nodeName => {
            if (!state.gateways.includes(nodeName)) return;
            const existing = state.gatewayNodes.find(gwn => gwn.node === nodeName);
            if (existing) {
                existing.don = workflowDONName;
            } else {
                state.gatewayNodes.push({ node: nodeName, don: workflowDONName });
            }
        });
    };

    function donColumnHTML(don, idx) {
        const donType = don.donTypes[0] || 'workflow';
        const nodesHTML = don.nodes.length === 0
            ? '<div class="don-nodes-empty">No chart nodes matched this DON\'s name</div>'
            : don.nodes.map(n => {
                const nodeInfo = state.nodes.find(s => s.name === n);
                return nodeInDonHTML(nodeInfo || { name: n, nodeType: 'standard' });
            }).join('');

        const capsHTML = don.capabilities.map(c => `
        <span class="cap-chip active" data-cap="${c}" data-don="${idx}">
            ${c}
            <span class="remove" onclick="removeCapability(${idx}, '${c}')">&times;</span>
        </span>
    `).join('');

        return `
        <div class="don-column" data-don-idx="${idx}">
            <div class="don-header">
                <div class="flex items-center gap-2 flex-1">
                    <input type="text" value="${don.name}"
                        class="text-sm font-semibold text-gray-800 border-none bg-transparent focus:ring-0 flex-1"
                        onchange="updateDONName(${idx}, this.value)">
                    <select class="don-type-select" onchange="updateDONType(${idx}, this.value)">
                        <option value="workflow" ${donType === 'workflow' ? 'selected' : ''}>workflow</option>
                        <option value="capabilities" ${donType === 'capabilities' ? 'selected' : ''}>capabilities</option>
                        <option value="gateway" ${donType === 'gateway' ? 'selected' : ''}>gateway</option>
                        <option value="bootstrap" ${donType === 'bootstrap' ? 'selected' : ''}>bootstrap</option>
                    </select>
                </div>
                <button onclick="removeDON(${idx})" class="text-gray-400 hover:text-red-500 text-sm ml-2">&times;</button>
            </div>
            <div class="don-body">
                <div class="don-nodes">
                    ${nodesHTML}
                </div>
                <div class="flex flex-wrap gap-1 mt-2">
                    ${capsHTML}
                    ${don.capabilities.length === 0 ? '<span class="text-xs text-gray-400">No capabilities selected</span>' : ''}
                </div>
                <div class="mt-2 flex items-center gap-2">
                    <label class="text-xs text-gray-500">
                        <input type="checkbox" ${don.exposesRemoteCapabilities ? 'checked' : ''}
                            onchange="toggleRemoteCaps(${idx}, this.checked)"
                            class="mr-1">
                        Exposes remote capabilities
                    </label>
                </div>
            </div>
        </div>
    `;
    }

    function nodeInDonHTML(node) {
        const role = node.nodeType || 'standard';
        const badgeClass = role === 'boot' ? 'badge-boot' : role === 'gateway' ? 'badge-gateway' : 'badge-standard';
        return `
        <div class="node-in-don" data-node="${node.name}">
            <div class="flex items-center gap-2">
                <span class="text-sm font-medium text-gray-800">${node.name}</span>
                <span class="badge ${badgeClass}">${role}</span>
            </div>
        </div>
    `;
    }

    // --- Render: Capability Catalog ---
    function renderCapabilityCatalog() {
        const container = document.getElementById('capabilityCatalog');
        if (!state.selectedDON && state.dons.length > 0) {
            state.selectedDON = 0;
        }

        if (state.dons.length === 0) {
            container.innerHTML = '<div class="text-xs text-gray-400 text-center py-4">Add a DON first</div>';
            return;
        }

        // DON selector
        let html = '<div class="mb-2"><select id="capDonSelect" onchange="selectDONForCaps(this.value)" class="don-type-select w-full">';
        state.dons.forEach((don, idx) => {
            html += `<option value="${idx}" ${state.selectedDON == idx ? 'selected' : ''}>${don.name}</option>`;
        });
        html += '</select></div>';

        const don = state.dons[state.selectedDON];
        if (!don) {
            container.innerHTML = html + '<div class="text-xs text-gray-400">Select a DON</div>';
            return;
        }

        // Build a set of active capability base names (strip chain suffix for comparison)
        const activeBases = new Set(don.capabilities.map(c => {
            for (const base of ['evm', 'solana', 'aptos']) {
                if (c.startsWith(base + '-')) return base;
            }
            return c;
        }));

        // For chain-scoped, also list which variants are active
        function chainVariants(base) {
            return don.capabilities.filter(c => c.startsWith(base + '-'));
        }

        html += '<div class="space-y-1">';
        state.capabilities.forEach(cap => {
            const isActive = activeBases.has(cap.name);
            const chainScoped = cap.chainScoped === true;
            const variants = chainScoped ? chainVariants(cap.name) : [];
            const hasVariants = variants.length > 0;
            const checkedIcon = (isActive || hasVariants)
                ? '<svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20"><path d="M16.7 5.3a1 1 0 010 1.4l-7.5 7.5a1 1 0 01-1.4 0L3.3 9.7a1 1 0 011.4-1.4l3.2 3.2 6.8-6.8a1 1 0 011.4 0z"/></svg>'
                : '';

            const variantList = hasVariants
                ? `<div class="flex flex-wrap gap-1 mt-1">${variants.map(v => `
                <span class="cap-chip active text-[10px]">${v}</span>
            `).join('')}</div>`
                : '';

            const labelSuffix = chainScoped
                ? ' <span class="text-[9px] text-purple-500">(chain)</span>'
                : '';

            html += `
            <div class="cap-catalog-item" onclick="toggleCapability(${state.selectedDON}, '${cap.name}', ${chainScoped})">
                <div class="flex items-center justify-between w-full">
                    <div class="flex-1">
                        <div class="text-xs font-medium text-gray-800">${cap.label}${labelSuffix}</div>
                        <div class="text-[10px] text-gray-400">${chainScoped ? 'Click to add chain' : cap.description}</div>
                        ${variantList}
                    </div>
                    <div class="w-4 h-4 rounded ${isActive || hasVariants ? 'bg-purple-600' : 'border border-gray-300'} flex items-center justify-center flex-shrink-0">
                        ${checkedIcon}
                    </div>
                </div>
            </div>
        `;
        });
        html += '</div>';

        container.innerHTML = html;
    }

    // --- Render: Config tab ---
    function renderConfig() {
        document.getElementById('cfgChartValues').value = state.infra.chartValues || '';
        document.getElementById('cfgNamespace').value = state.infra.namespace || '';
        document.getElementById('cfgJDGrpc').value = state.jd.grpc || '';
        document.getElementById('cfgJDDomain').value = state.jd.domain || '';
        document.getElementById('cfgJDEnv').value = state.jd.environment || '';
        document.getElementById('cfgJDUseTLS').checked = state.jd.useTLS !== false;
    }

    // --- Render: Cap Configs tab ---
    function renderCapConfigs() {
        const container = document.getElementById('capConfigsContent');
        const entries = Object.entries(state.capabilityConfigs).sort((a, b) => a[0].localeCompare(b[0]));

        if (entries.length === 0) {
            container.innerHTML = '<div class="text-xs text-gray-400 text-center py-4">No capability configs — defaults will be loaded on save</div>';
            return;
        }

        container.innerHTML = entries.map(([name, cc]) => capConfigCardHTML(name, cc)).join('');
    }

    function capConfigCardHTML(name, cc) {
        const values = cc.Values || {};
        const valueEntries = Object.entries(values);

        function formatValue(v) {
            if (v === null || v === undefined) return '';
            if (typeof v === 'object') return JSON.stringify(v, null, 2);
            return String(v);
        }

        function isObject(v) {
            return v !== null && typeof v === 'object';
        }

        const valuesHTML = valueEntries.length === 0
            ? '<div class="text-[10px] text-gray-400 italic">No values configured</div>'
            : valueEntries.map(([k, v]) => {
                const obj = isObject(v);
                const valStr = formatValue(v);
                const valId = `capval-${name}-${k}`;
                return `
                <div class="mb-2" data-cap-val-row="${name}|${k}">
                    <div class="flex items-center gap-2 mb-1">
                        <input type="text" value="${k}"
                            class="cap-config-key flex-1 text-xs px-2 py-1 border border-gray-200 rounded"
                            data-cap-config="${name}" data-cap-oldkey="${k}" data-field="key"
                            placeholder="key">
                        <button onclick="removeCapConfigValue('${name}', '${k}')"
                            class="text-gray-300 hover:text-red-500 text-sm">&times;</button>
                    </div>
                    <textarea rows="${obj ? 4 : 1}"
                        class="cap-config-val w-full text-xs px-2 py-1 border border-gray-200 rounded font-mono resize-y"
                        data-cap-config="${name}" data-cap-key="${k}" data-field="value"
                        placeholder="value (JSON for objects)">${valStr}</textarea>
                </div>
            `;
            }).join('');

        return `
        <div class="border border-gray-200 rounded-lg p-3" data-cap-config-card="${name}">
            <div class="flex items-center justify-between mb-2">
                <input type="text" value="${name}"
                    class="text-sm font-semibold text-gray-800 border-none bg-transparent focus:ring-0"
                    data-cap-config="${name}" data-field="name" data-oldname="${name}">
                <button onclick="removeCapConfig('${name}')"
                    class="text-gray-300 hover:text-red-500 text-sm">&times;</button>
            </div>
            <div class="mb-2">
                <label class="text-[10px] text-gray-400 uppercase tracking-wide">Binary Name</label>
                <input type="text" value="${cc.BinaryName || ''}"
                    class="w-full text-xs px-2 py-1 border border-gray-200 rounded font-mono"
                    data-cap-config="${name}" data-field="binary">
            </div>
            <div>
                <div class="flex items-center justify-between mb-1">
                    <label class="text-[10px] text-gray-400 uppercase tracking-wide">Values</label>
                    <button onclick="addCapConfigValue('${name}')"
                        class="text-[10px] text-blue-500 hover:text-blue-700">+ Add value</button>
                </div>
                ${valuesHTML}
            </div>
        </div>
    `;
    }

    // --- Cap Config actions ---

    // syncCapConfigsFromDOM reads all cap config inputs from the DOM and writes
    // them into state.capabilityConfigs. Called before save/preview to ensure
    // unsaved input changes (which haven't triggered onchange yet) are captured.
    function syncCapConfigsFromDOM() {
        const cards = document.querySelectorAll('[data-cap-config-card]');
        // If the Cap Configs tab hasn't been rendered yet (no cards in DOM),
        // don't touch state.capabilityConfigs — keep whatever was loaded from
        // the API (defaults or saved desired state).
        if (cards.length === 0) return;

        const newConfigs = {};
        cards.forEach(card => {
            const nameInput = card.querySelector('[data-field="name"]');
            if (!nameInput) return;
            const name = nameInput.value.trim();
            if (!name) return;

            const binaryInput = card.querySelector('[data-field="binary"]');
            const binary = binaryInput ? binaryInput.value : '';

            const values = {};
            const valRows = card.querySelectorAll('[data-cap-val-row]');
            valRows.forEach(row => {
                const keyInput = row.querySelector('[data-field="key"]');
                const valInput = row.querySelector('[data-field="value"]');
                if (!keyInput || !valInput) return;
                const key = keyInput.value.trim();
                if (!key) return;
                const rawVal = valInput.value.trim();
                // Try to parse as JSON (for nested objects/numbers), otherwise keep as string
                let parsed = rawVal;
                if (rawVal) {
                    try { parsed = JSON.parse(rawVal); } catch (e) { }
                }
                values[key] = parsed;
            });

            newConfigs[name] = { BinaryName: binary, Values: values };
        });
        state.capabilityConfigs = newConfigs;
    }

    window.removeCapConfigValue = function (name, key) {
        if (!state.capabilityConfigs[name]) return;
        // Sync first to capture any edits, then remove
        syncCapConfigsFromDOM();
        if (state.capabilityConfigs[name] && state.capabilityConfigs[name].Values) {
            delete state.capabilityConfigs[name].Values[key];
        }
        renderCapConfigs();
    };

    window.addCapConfigValue = function (name) {
        // Sync first to capture any edits
        syncCapConfigsFromDOM();
        if (!state.capabilityConfigs[name]) return;
        if (!state.capabilityConfigs[name].Values) state.capabilityConfigs[name].Values = {};
        const newKey = 'new_key_' + Date.now();
        state.capabilityConfigs[name].Values[newKey] = '';
        renderCapConfigs();
    };

    window.removeCapConfig = function (name) {
        syncCapConfigsFromDOM();
        delete state.capabilityConfigs[name];
        renderCapConfigs();
    };

    window.addCapConfig = function () {
        const name = prompt('Enter capability config name (e.g. "evm", "cron"):');
        if (!name) return;
        syncCapConfigsFromDOM();
        if (state.capabilityConfigs[name]) {
            toast(`${name} already exists`, 'error');
            return;
        }
        state.capabilityConfigs[name] = { BinaryName: '', Values: {} };
        renderCapConfigs();
    };

    // --- Render: Status tab ---
    async function loadState() {
        const container = document.getElementById('statusContent');
        try {
            const s = await apiGet('/api/state');
            if (!s.hasState) {
                container.innerHTML = `
                <div class="status-card">
                    <h3>Reconcile State</h3>
                    <p class="text-sm text-gray-500">No state file found. Run <code class="bg-gray-100 px-1 rounded">reconciler apply</code> to start.</p>
                </div>`;
                return;
            }

            const phaseClass = s.phase === 'done' ? 'phase-done' :
                s.phase === 'toml' ? 'phase-toml' :
                    s.phase === 'on-chain' ? 'phase-on-chain' : 'phase-none';

            let addrsHTML = '';
            if (s.addresses && s.addresses.length > 0) {
                addrsHTML = s.addresses.map(a => `
                <div class="status-row">
                    <span class="status-label">${a.type}</span>
                    <span class="status-value font-mono text-xs">${a.address}</span>
                </div>
            `).join('');
            } else {
                addrsHTML = '<div class="text-sm text-gray-400">No contracts deployed yet</div>';
            }

            let donIDsHTML = '';
            if (s.donIds && Object.keys(s.donIds).length > 0) {
                donIDsHTML = Object.entries(s.donIds).map(([name, id]) => `
                <div class="status-row">
                    <span class="status-label">${name}</span>
                    <span class="status-value">DON ID: ${id}</span>
                </div>
            `).join('');
            } else {
                donIDsHTML = '<div class="text-sm text-gray-400">No DONs registered on-chain yet</div>';
            }

            container.innerHTML = `
            <div class="grid grid-cols-2 gap-4">
                <div class="status-card">
                    <h3>Reconcile Phase</h3>
                    <div class="text-center py-4">
                        <span class="phase-badge ${phaseClass}">${s.phase || 'not started'}</span>
                    </div>
                </div>
                <div class="status-card">
                    <h3>Deployed Contracts</h3>
                    ${addrsHTML}
                </div>
                <div class="status-card">
                    <h3>On-Chain DON IDs</h3>
                    ${donIDsHTML}
                </div>
                <div class="status-card">
                    <h3>Nodes</h3>
                    ${state.nodes.map(n => `
                        <div class="status-row">
                            <span class="status-label">${n.name}</span>
                            <span class="status-value">${n.nodeType}</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
        } catch (e) {
            container.innerHTML = `<div class="status-card"><p class="text-sm text-red-500">${e.message}</p></div>`;
        }
    }

    // --- Actions ---

    window.updateDONName = function (idx, name) {
        state.dons[idx].name = name;
        renderCapabilityCatalog();
    };

    window.updateDONType = function (idx, type) {
        state.dons[idx].donTypes = [type];
    };

    window.removeDON = function (idx) {
        state.dons.splice(idx, 1);
        if (state.selectedDON >= state.dons.length) state.selectedDON = Math.max(0, state.dons.length - 1);
        renderDONs();
    };

    // donTypeForNodes derives a DON's type from the roles of its member nodes —
    // a node-set with a gateway node is a gateway DON, one with a boot node is a
    // bootstrap DON, otherwise it's a workflow DON.
    function donTypeForNodes(nodes) {
        if (nodes.some(n => n.nodeType === 'gateway')) return 'gateway';
        if (nodes.some(n => n.nodeType === 'boot')) return 'bootstrap';
        return 'workflow';
    }

    // autoDetectDONs groups chart nodes by their Griddle-registered don-name
    // label and creates one DON per node-set not already present. This is the
    // only way DONs get created — there's no manual "add DON" action, since a
    // DON's name must exactly match what Griddle registered with JD, and its
    // membership is always read from the chart (never edited here).
    window.autoDetectDONs = function () {
        const existing = new Set(state.dons.map(d => d.name));

        const groups = new Map(); // chartDonName -> node objects
        state.nodes.forEach(n => {
            if (!n.chartDonName || existing.has(n.chartDonName)) return;
            if (!groups.has(n.chartDonName)) groups.set(n.chartDonName, []);
            groups.get(n.chartDonName).push(n);
        });

        if (groups.size === 0) return;

        groups.forEach((nodes, donName) => {
            state.dons.push({
                name: donName,
                donTypes: [donTypeForNodes(nodes)],
                capabilities: [],
                nodes: nodes.map(n => n.name),
                bootstrapNode: '',
                exposesRemoteCapabilities: false,
                registryBasedLaunchAllowlist: [],
                capabilityConfigs: {},
            });
        });

        renderDONs();
        toast(`Added ${groups.size} DON(s) from chart node-sets`, 'success');
    };

    window.toggleCapability = function (donIdx, capName, chainScoped) {
        const don = state.dons[donIdx];

        if (chainScoped) {
            // Only chains declared in the Chains section can be picked — this is
            // what keeps a capability from ever pointing at an undeclared chain
            // (the root cause of the evm-1337-vs-registry-chain bug this UI exists
            // to prevent).
            if (state.chains.length === 0) {
                toast('No chains declared yet — add one in the Chains section first', 'error');
                return;
            }
            const options = state.chains.map(c => String(c.chainId));
            const chainId = prompt(`Enter chain ID for ${capName} capability (declared: ${options.join(', ')}):`);
            if (!chainId || isNaN(parseInt(chainId))) {
                if (chainId !== null) toast('Invalid chain ID', 'error');
                return;
            }
            if (!options.includes(String(parseInt(chainId)))) {
                toast(`Chain ${chainId} is not declared — add it in the Chains section first`, 'error');
                return;
            }
            const fullCap = capName + '-' + chainId;
            if (don.capabilities.includes(fullCap)) {
                toast(`${fullCap} already added`, 'error');
                return;
            }
            don.capabilities.push(fullCap);
        } else {
            // Non-chain-scoped: toggle on/off
            const idx = don.capabilities.indexOf(capName);
            if (idx >= 0) {
                don.capabilities.splice(idx, 1);
            } else {
                don.capabilities.push(capName);
            }
        }
        renderDONs();
    };

    window.removeCapability = function (donIdx, capName) {
        const don = state.dons[donIdx];
        don.capabilities = don.capabilities.filter(c => c !== capName);
        renderDONs();
    };

    window.toggleRemoteCaps = function (idx, checked) {
        state.dons[idx].exposesRemoteCapabilities = checked;
    };

    window.selectDONForCaps = function (idx) {
        state.selectedDON = parseInt(idx);
        renderCapabilityCatalog();
    };

    // --- Save ---
    async function save() {
        syncCapConfigsFromDOM();
        syncConfigInputs();
        const payload = { infra: state.infra, jd: state.jd, chains: state.chains, dons: state.dons, gatewayNodes: state.gatewayNodes, capabilityConfigs: state.capabilityConfigs };
        try {
            await apiPost('/api/desired', payload);
            toast('Desired state saved', 'success');
        } catch (e) {
            toast('Save failed: ' + e.message, 'error');
        }
    }

    // --- Preview TOML ---
    async function previewTOML() {
        syncCapConfigsFromDOM();
        syncConfigInputs();
        const payload = { infra: state.infra, jd: state.jd, chains: state.chains, dons: state.dons, gatewayNodes: state.gatewayNodes, capabilityConfigs: state.capabilityConfigs };
        try {
            const resp = await fetch('/api/preview-toml', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });
            const text = await resp.text();
            document.getElementById('previewContent').textContent = text;
            document.getElementById('previewModal').classList.remove('hidden');
            document.getElementById('previewModal').classList.add('flex');
        } catch (e) {
            toast('Preview failed: ' + e.message, 'error');
        }
    }

    function syncConfigInputs() {
        state.infra.chartValues = document.getElementById('cfgChartValues').value;
        state.infra.namespace = document.getElementById('cfgNamespace').value;
        state.jd.grpc = document.getElementById('cfgJDGrpc').value;
        state.jd.domain = document.getElementById('cfgJDDomain').value;
        state.jd.environment = document.getElementById('cfgJDEnv').value;
        state.jd.useTLS = document.getElementById('cfgJDUseTLS').checked;
    }

    // --- Check JD Connectivity ---
    async function checkJD() {
        syncConfigInputs();

        const btn = document.getElementById('btnCheckJD');
        const resultEl = document.getElementById('jdCheckResult');
        const nodeResultsEl = document.getElementById('jdNodeResults');
        const nodeListEl = document.getElementById('jdNodeList');

        btn.disabled = true;
        btn.textContent = 'Checking...';
        nodeResultsEl.classList.add('hidden');

        // Collect node names — use DON-assigned nodes if any, otherwise all discovered nodes
        const nodeNames = [];
        const seen = new Set();
        state.dons.forEach(don => don.nodes.forEach(n => {
            if (!seen.has(n)) { nodeNames.push(n); seen.add(n); }
        }));
        if (nodeNames.length === 0) {
            state.nodes.forEach(n => {
                if (!seen.has(n.name)) { nodeNames.push(n.name); seen.add(n.name); }
            });
        }

        // Phase 1: JD connectivity
        resultEl.innerHTML = '<span class="text-gray-500"><span class="jd-spinner"></span> Connecting to JD...</span>';

        try {
            const resp = await apiPost('/api/jd/check', {
                grpc: state.jd.grpc,
                useTLS: state.jd.useTLS,
                namespace: state.infra.namespace,
                kubeconfig: '',
                nodeNames: nodeNames,
            });

            if (!resp.connected) {
                let html = `<span class="text-red-600 font-medium">✗ JD connection failed</span>`;
                html += `<div class="mt-1 text-xs text-red-500">${resp.error}</div>`;
                if (resp.k8sErrors && resp.k8sErrors.length > 0) {
                    html += `<div class="mt-1 text-xs text-gray-500">K8s: ${resp.k8sErrors.join('; ')}</div>`;
                }
                resultEl.innerHTML = html;
                return;
            }

            // JD connected — show progress for node validation
            resultEl.innerHTML = '<span class="text-green-600 font-medium">✓ JD connected</span> <span class="text-gray-500"><span class="jd-spinner"></span> Validating nodes...</span>';

            // Small delay so the user sees the progressive update
            await new Promise(r => setTimeout(r, 300));

            if (resp.error && (!resp.nodes || resp.nodes.length === 0)) {
                let html = `<span class="text-green-600 font-medium">✓ JD connected</span>`;
                html += `<div class="mt-1 text-xs text-yellow-600">${resp.error}</div>`;
                if (resp.k8sErrors && resp.k8sErrors.length > 0) {
                    html += `<div class="mt-1 text-xs text-gray-500">K8s: ${resp.k8sErrors.join('; ')}</div>`;
                }
                resultEl.innerHTML = html;
                return;
            }

            // Show per-node results progressively
            const nodes = resp.nodes || [];
            const found = nodes.filter(n => n.found).length;
            const total = nodes.length;
            const allOk = nodes.every(n => n.found && n.isConnected);

            let summary = '';
            if (allOk) {
                summary = `<span class="text-green-600 font-medium">✓ All ${total} nodes found and connected</span>`;
            } else {
                summary = `<span class="text-yellow-600 font-medium">⚠ ${found}/${total} nodes validated</span>`;
            }
            if (resp.k8sErrors && resp.k8sErrors.length > 0) {
                summary += `<div class="mt-1 text-xs text-gray-500">K8s: ${resp.k8sErrors.join('; ')}</div>`;
            }
            resultEl.innerHTML = summary;

            // Render node results one by one with a small delay
            nodeListEl.innerHTML = '';
            nodeResultsEl.classList.remove('hidden');

            for (let i = 0; i < nodes.length; i++) {
                const n = nodes[i];
                const icon = n.found ? (n.isConnected ? '✓' : '⚠') : '✗';
                const color = n.found ? (n.isConnected ? 'text-green-600' : 'text-yellow-600') : 'text-red-600';
                const csaShort = n.csaKey ? n.csaKey.substring(0, 12) + '...' : 'N/A';
                const row = document.createElement('div');
                row.className = 'flex items-center gap-3 px-3 py-2 border border-gray-200 rounded-lg';
                row.style.opacity = '0';
                row.style.transition = 'opacity 0.3s';
                row.innerHTML = `
                <span class="${color} text-lg font-bold">${icon}</span>
                <div class="flex-1">
                    <div class="text-sm font-medium text-gray-800">${n.nodeName}</div>
                    <div class="text-xs text-gray-400">CSA: ${csaShort}</div>
                </div>
                <div class="text-right">
                    ${n.found ? `<div class="text-xs text-gray-600">JD: ${n.jdName || n.jdId}</div>` : '<div class="text-xs text-red-500">not in JD</div>'}
                    ${n.error ? `<div class="text-xs text-yellow-600">${n.error}</div>` : ''}
                </div>
            `;
                nodeListEl.appendChild(row);
                // Fade in
                requestAnimationFrame(() => { row.style.opacity = '1'; });
                if (i < nodes.length - 1) {
                    await new Promise(r => setTimeout(r, 200));
                }
            }

        } catch (e) {
            resultEl.innerHTML = `<span class="text-red-600 font-medium">✗ Error: ${e.message}</span>`;
        } finally {
            btn.disabled = false;
            btn.textContent = 'Check JD Connectivity';
        }
    }

    // --- Init ---
    function init() {
        // Tab listeners
        document.querySelectorAll('[data-tab]').forEach(el => {
            el.addEventListener('click', () => switchTab(el.dataset.tab));
        });

        // Button listeners
        document.getElementById('btnSave').addEventListener('click', save);
        document.getElementById('btnPreview').addEventListener('click', previewTOML);
        document.getElementById('btnAddCapConfig').addEventListener('click', window.addCapConfig);
        document.getElementById('btnCheckJD').addEventListener('click', checkJD);
        document.getElementById('btnLoadChains').addEventListener('click', () => loadDiscoveredChains(false));
        document.getElementById('btnClosePreview').addEventListener('click', () => {
            document.getElementById('previewModal').classList.add('hidden');
            document.getElementById('previewModal').classList.remove('flex');
        });

        // Load data, then auto-detect DONs from the chart's don-name labels —
        // the only way DONs get populated (no manual "add DON" action). Chains
        // are always preloaded from the chart's node config on startup — there
        // is no manual "add chain" action, since the chart is the source of truth.
        loadNodes().then(async () => {
            await loadDesired();
            loadCapabilities();
            window.autoDetectDONs();
            loadDiscoveredChains(true);
        });
    }

    document.addEventListener('DOMContentLoaded', init);
})();
