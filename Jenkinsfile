pipeline {
    agent any

    options {
        buildDiscarder(logRotator(
            numToKeepStr: '10',
            daysToKeepStr: '7',
            artifactNumToKeepStr: '5'
        ))
        timeout(time: 30, unit: 'MINUTES')
    }

    environment {
        // 镜像配置
        IMAGE_NAME = 'video-backend'
        APP_NAME = 'video-backend'  // 容器名称

        // 部署配置
        DEPLOY_HOST = 'localhost'          // 本地部署
        DEPLOY_USER = ''                   // 远程部署时的SSH用户名
        SSH_CREDENTIALS_ID = ''           // 远程部署时的SSH凭据ID

        // 配置文件凭据ID
        CONFIG_CREDENTIAL_ID = 'video-backend-config'
    }

    stages {
        stage('检出代码') {
            steps {
                deleteDir() // 删除工作空间
                checkout scm // 检出代码
            }
        }
        
        stage('测试') {
            steps {
                // Go 的测试命令，运行项目中的测试文件（*_test.go），-v 参数用于显示详细输出
                sh 'go test -v ./... || true'
            }
            post {
                always {
                    // 测试后清理配置文件，避免污染 git
                    sh 'rm -f config.yaml || true'
                    sh 'ls -la ./'
                }
            }
        }

        stage('准备配置文件') {
            steps {
                withCredentials([file(credentialsId: "${CONFIG_CREDENTIAL_ID}", variable: 'CONFIG_FILE')]) {
                    sh 'cp "$CONFIG_FILE" ./config.yaml'
                }
            }
        }

        stage('构建Docker镜像') {
            steps {
                sh "docker build -t ${IMAGE_NAME}:latest ."
            }
        }

        stage('部署') {
            steps {
                script {
                    def isLocal = !"${DEPLOY_HOST}" || "${DEPLOY_HOST}" == 'localhost' || "${DEPLOY_HOST}" == '127.0.0.1'
                    
                    if (isLocal) {
                        sh """
                            export IMAGE=${IMAGE_NAME}:latest
                            export APP_NAME=${APP_NAME}
                            docker compose down -v || true
                            docker compose up -d
                            docker compose ps
                        """
                    } else {
                        sshagent(["${SSH_CREDENTIALS_ID}"]) {
                            sh """
                                echo "远程部署待实现"
                                exit 1
                            """
                        }
                    }
                }
            }
        }

        stage('服务状态') {
            steps {
                sh """
                    sleep 5
                    for i in 1 2 3 4 5 6 7 8 9 10; do
                        curl -f -s http://host.docker.internal:8888/health > /dev/null && echo "✅ 服务健康检查通过" && exit 0
                        sleep 2
                    done
                    echo "❌ 服务健康检查失败"
                    docker compose logs --tail=20
                    exit 1
                """
            }
        }
    }

    post {
        always {
            sh 'rm -f config.yaml || true'
            cleanWs()
        }
        success {
            echo "✅ 构建成功！镜像: ${IMAGE_NAME}:latest"
        }
        failure {
            echo "❌ 构建失败！"
        }
    }
}