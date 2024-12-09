pipeline {
    agent any
    stages {
        stage('Build and Test') {
            steps {
                script {
                    try {
                        // Construye la imagen, incluyendo pruebas en el Dockerfile
                        sh 'docker build -t api-image2 .'
                    } catch (Exception err) {
                        error("Error en la construcción o en las pruebas.")
                    }
                }
            }
        }
    }
    post {
        success {
            echo 'Build y pruebas exitosas.'
        }
        failure {
            echo 'Error en la construcción o en las pruebas.'
        }
    }
}
