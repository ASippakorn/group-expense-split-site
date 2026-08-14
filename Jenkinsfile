pipeline {
  agent any

  stages {
    stage('API Tests') {
      steps {
        dir('apps/api') {
          sh 'go test ./...'
        }
      }
    }

    stage('Web Tests') {
      steps {
        sh 'corepack enable'
        sh 'pnpm install --frozen-lockfile'
        sh 'pnpm --filter @splitr/web test'
      }
    }

    stage('OpenAPI Lint') {
      steps {
        sh 'pnpm openapi:lint'
      }
    }
  }
}
