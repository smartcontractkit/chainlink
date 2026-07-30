./run.sh runs a workflow building the capabilities locally (use rtinianov_cre_psb_ms_1_p branch), it'll continually print the result of the workflow until ctrl+C is pressed. The local environment will be kill at that point.


By default, all nodes proxy Rage to the PSB (to allow further capability development using it).
Use ./two-proxy.sh and ./four-proxy.sh to change the settings. two-proxy.sh sets nodes 2 and 3 to proxy in the workflow and capability DONs.


