# shellcheck shell=sh

Describe 'Scheduled'
  It 'Trigger on schedule when source changed'
    export SRC_CHANGED=true
    export GITHUB_EVENT_NAME=schedule 
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=true'
  End

  It 'Trigger on schedule when source unchanged'
    export SRC_CHANGED=false
    export GITHUB_EVENT_NAME=schedule 
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=true'
  End
End


Describe 'Pull Request'
  It 'Trigger on pull_request for non release branches when source changed'
    export SRC_CHANGED=true
    export GITHUB_EVENT_NAME=pull_request
    export GITHUB_BASE_REF=develop
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=true'
  End

  It 'No trigger on pull_request for non release branches when source unchanged'
    export SRC_CHANGED=false
    export GITHUB_EVENT_NAME=pull_request
    export GITHUB_BASE_REF=develop
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=false'
  End

  It 'No trigger on pull_request for release branches when source changed'
    export SRC_CHANGED=true
    export GITHUB_EVENT_NAME=pull_request
    export GITHUB_BASE_REF=release/1.2.3
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=false'
  End
End

Describe 'Push'
  It 'Trigger on push to the staging branch when source changed'
    export SRC_CHANGED=true
    export GITHUB_EVENT_NAME=push
    export GITHUB_REF=staging
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=true'
  End

  It 'No trigger on push to the staging branch when source unchanged'
    export SRC_CHANGED=false
    export GITHUB_EVENT_NAME=push
    export GITHUB_REF=staging
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=false'
  End

  It 'Trigger on push to the trying branch when source changed'
    export SRC_CHANGED=true
    export GITHUB_EVENT_NAME=push
    export GITHUB_REF=trying
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=true'
  End

  It 'Trigger on push to the rollup branch when source changed'
    export SRC_CHANGED=true
    export GITHUB_EVENT_NAME=push
    export GITHUB_REF=rollup
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=true'
  End

  It 'No trigger on push to the develop branch when source changed'
    export SRC_CHANGED=true
    export GITHUB_EVENT_NAME=push
    export GITHUB_REF=develop
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=false'
  End
End

Describe 'Misc'
  It 'No trigger on invalid event name when source changed'
    export SRC_CHANGED=true
    export GITHUB_EVENT_NAME=invalid_event_name
    When run ./ontriggerlint.sh
    The status should eq 0
    The stdout should equal 'on_trigger=false'
  End
End
