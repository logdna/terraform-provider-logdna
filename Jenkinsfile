library 'magic-butler-catalogue'
def PROJECT_NAME = 'terraform-provider-logdna'
def CURRENT_BRANCH = currentBranch()
def TRIGGER_PATTERN = ".*@logdnabot.*"

pipeline {
  agent none

  options {
    timestamps()
    ansiColor 'xterm'
  }

  triggers {
    issueCommentTrigger(TRIGGER_PATTERN)
  }

  stages {
    stage('Validate PR Source') {
      when {
        expression { env.CHANGE_FORK }
        not {
          triggeredBy 'issueCommentCause'
        }
      }
      steps {
        error("A maintainer needs to approve this PR for CI by commenting")
      }
    }

    stage('Test Suite') {
      matrix {
        axes {
          axis {
              name 'GO_VERSION'
              values '1.14', '1.16'
          }
        }

        agent {
          node {
            label 'ec2-fleet'
            customWorkspace "${PROJECT_NAME}-${BUILD_NUMBER}-${GO_VERSION}"
          }
        }

        environment {
          SERVICE_KEY = credentials('terraform-provider-servicekey')
          COVERALLS_TOKEN = credentials('terraform-provider-coveralls')
        }

        stages {
          stage('Test') {
            steps {
              script {
                compose.up(
                  PROJECT_NAME
                , ['compose/test.yml']
                , ['build': true]
                )
              }
            }

            post {
              always {
                script {
                  compose.down(
                    ['compose/test.yml']
                  , [('remove-orphans'): true, ('volumes'): true, ('rmi'): 'local']
                  )
                }
                publishHTML target: [
                  allowMissing: false,
                  alwaysLinkToLastBuild: false,
                  keepAll: true,
                  reportDir: 'coverage/',
                  reportFiles: '*.html',
                  reportName: "Code Coverage for Go ${GO_VERSION}"
                ]
              }
            }
          }
        }
      }
    }
  }
}