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
        error('A maintainer needs to approve this PR for CI by commenting')
      }
    }

    stage('Test') {
      agent {
        node {
          label 'ec2-fleet'
          customWorkspace "${PROJECT_NAME}-${BUILD_NUMBER}"
        }
      }

      environment {
        GIT_BRANCH = "${CURRENT_BRANCH}"
      }

      steps {
        script {
          withCredentials([
            string(credentialsId: 'logdna-gpg-key', variable: 'GPG_KEY'),
            string(credentialsId: 'terraform-provider-servicekey', variable: 'SERVICE_KEY'),
            string(credentialsId: 'terraform-provider-coveralls', variable: 'COVERALLS_TOKEN')
          ]) {
            sh '''
              set +x
              echo "$GPG_KEY" > gpgkey.asc              
              make postcov
              make test-release
            '''
          }
        }
      }

      post {
        always {
          sh 'rm -f gpgkey.asc'
          publishHTML target: [
            allowMissing: false,
            alwaysLinkToLastBuild: false,
            keepAll: true,
            reportDir: 'coverage/',
            reportFiles: '*.html',
            reportName: 'Code Coverage'
          ]
          archiveArtifacts 'dist/*'
        }
      }
    }

    stage('Release') {
      when {
        beforeAgent true
        buildingTag()
      }

      agent {
        node {
          label 'ec2-fleet'
          customWorkspace "${PROJECT_NAME}-${BUILD_NUMBER}"
        }
      }

      steps {
        script {
          withCredentials([
            string(credentialsId: 'logdna-gpg-key', variable: 'GPG_KEY'),
            string(credentialsId: 'github-api-token', variable: 'GITHUB_TOKEN')
          ]) {
            sh '''
              set +x
              echo "$GPG_KEY" > gpgkey.asc              
              make release
            '''
          }
        }
      }

      post {
        always {
          sh 'rm -f gpgkey.asc'
        }
      }
    }
  }
}
