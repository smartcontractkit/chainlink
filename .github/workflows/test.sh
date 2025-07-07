          if [[ "${GITHUB_EVENT_NAME}" == 'pull_request' ]]; then
            # Run Core E2E tests on PRs only if the label "run-e2e-tests" is present, and there are relevant changes
            echo "event-shorthand=PR" | tee -a "$GITHUB_OUTPUT"
            echo "workflow-name=Run Core E2E Tests For PR" | tee -a "$GITHUB_OUTPUT"
            echo "test-trigger=PR E2E Core Tests" | tee -a "$GITHUB_OUTPUT"

            if [[ "${RUN_E2E_TESTS_LABEL_FOUND}" == 'true' && "${CONTAINS_CHANGES}" == 'true' ]]; then
              echo "should-run=true" | tee -a "$GITHUB_OUTPUT"
            fi

          elif [[ "${GITHUB_EVENT_NAME}" == 'merge_group' ]]; then
            # Run Core E2E tests in the merge queue, if there are relevant changes
            echo "event-shorthand=Merge Queue" | tee -a "$GITHUB_OUTPUT"
            echo "workflow-name=Run Core CRE Tests For Merge Queue" | tee -a "$GITHUB_OUTPUT"
            echo "test-trigger=Merge Queue CRE E2E Core Tests" | tee -a "$GITHUB_OUTPUT"
            echo "should-run=${CONTAINS_CHANGES}" | tee -a "$GITHUB_OUTPUT"

          elif [[ "${GITHUB_EVENT_NAME}" == 'workflow_dispatch' ]]; then
            # Always run Core E2E tests on workflow dispatch
            echo "event-shorthand=Workflow Dispatch" | tee -a "$GITHUB_OUTPUT"
            echo "workflow-name=Run Core E2E Tests For Workflow Dispatch" | tee -a "$GITHUB_OUTPUT"
            echo "test-trigger=Workflow Dispatch E2E Core Tests" | tee -a "$GITHUB_OUTPUT"
            echo "should-run=true" | tee -a "$GITHUB_OUTPUT"

          elif [[ "${GITHUB_EVENT_NAME}" == 'push' ]]; then
            # Run Core E2E tests on push events, only if there are relevant changes or if it's a tag push
            echo "event-shorthand=Push" | tee -a "$GITHUB_OUTPUT"
            echo "workflow-name=Run Core E2E Tests For Push" | tee -a "$GITHUB_OUTPUT"
            echo "test-trigger=Push E2E Core Tests" | tee -a "$GITHUB_OUTPUT"
            if [[ "${CONTAINS_CHANGES}" == 'true' || "${GITHUB_REF_TYPE}" == 'tag' ]]; then
              echo "should-run=true" | tee -a "$GITHUB_OUTPUT"
            fi

          else
            echo "event-shorthand=Unknown" | tee -a "$GITHUB_OUTPUT"
            echo "workflow-name=Run Core E2E Tests For Unknown Event" | tee -a "$GITHUB_OUTPUT"
            echo "test-trigger=Unknown E2E Core Tests" | tee -a "$GITHUB_OUTPUT"
            echo "should-run=false" | tee -a "$GITHUB_OUTPUT"
          fi
