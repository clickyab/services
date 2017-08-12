node {
    stage('Build') {
        checkout scm
        sh "./bin/ci-test.sh all"
    }
    stage('Lint') {
        checkout scm
        sh "./bin/ci-test.sh lint"
    }
    stage('Test') {
        checkout scm
        def REDIS_NAME = sh(script: 'cat /dev/urandom | tr -dc "a-zA-Z0-9" | fold -w 32 | head -n 1', returnStdout: true).trim()
        sh "docker rm -f $REDIS_NAME || true"
        sh "docker run -d --rm --name $REDIS_NAME -p 6379:6379 redis:alpine"
        sh "./bin/ci-test.sh test"
        sh "docker rm -f $REDIS_NAME"
    }
}
