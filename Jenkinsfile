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
        IMAGE_TAG = "${env.BUILD_NUMBER}"
        
        // 部署配置
        // 如果DEPLOY_HOST为空或localhost，则为本地部署；否则为远程部署
        DEPLOY_HOST = 'localhost'          // 部署服务器地址，空或localhost=本地部署，否则=远程部署
        DEPLOY_USER = ''                   // 远程部署时的SSH用户名（远程部署必填）
        SSH_CREDENTIALS_ID = ''           // 远程部署时的SSH凭据ID（远程部署必填）

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
        
        stage('运行测试') {
            steps {
                sh 'go test -v ./... || true'
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
                sh """
                    docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
                    docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${IMAGE_NAME}:latest
                """
            }
        }

        stage('部署') {
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                script {
                    def isLocal = !"${DEPLOY_HOST}" || "${DEPLOY_HOST}" == 'localhost' || "${DEPLOY_HOST}" == '127.0.0.1'
                    
                    withCredentials([file(credentialsId: "${CONFIG_CREDENTIAL_ID}", variable: 'CONFIG_FILE')]) {
                        if (isLocal) {
                            // 本地部署：直接使用当前工作目录
                            sh """
                                cp "\$CONFIG_FILE" ./config.yaml
                                IMAGE=${IMAGE_NAME}:latest docker-compose down || true
                                IMAGE=${IMAGE_NAME}:latest docker-compose up -d
                                docker-compose ps
                            """
                        } else {
                            // 远程部署：直接在远程服务器执行
                            sshagent(["${SSH_CREDENTIALS_ID}"]) {
                                sh """
                                    set -e
                                    IMAGE_TAR=${IMAGE_NAME}-${IMAGE_TAG}.tar.gz
                                    
                                    # 保存并传输镜像
                                    docker save ${IMAGE_NAME}:${IMAGE_TAG} ${IMAGE_NAME}:latest | gzip > "\${IMAGE_TAR}"
                                    scp "\${IMAGE_TAR}" ${DEPLOY_USER}@${DEPLOY_HOST}:/tmp/
                                    
                                    # 传输配置文件
                                    scp "\$CONFIG_FILE" ${DEPLOY_USER}@${DEPLOY_HOST}:/tmp/config.yaml
                                    scp docker-compose.yml ${DEPLOY_USER}@${DEPLOY_HOST}:/tmp/docker-compose.yml
                                    
                                    # 在远程服务器执行部署
                                    ssh ${DEPLOY_USER}@${DEPLOY_HOST} << EOF
                                        set -e
                                        cd /tmp
                                        docker load < /tmp/${IMAGE_NAME}-${IMAGE_TAG}.tar.gz
                                        rm -f /tmp/${IMAGE_NAME}-${IMAGE_TAG}.tar.gz
                                        IMAGE=${IMAGE_NAME}:latest docker-compose down || true
                                        IMAGE=${IMAGE_NAME}:latest docker-compose up -d
                                        docker-compose ps
                                    EOF
                                    
                                    rm -f "\${IMAGE_TAR}"
                                """
                            }
                        }
                    }
                }
            }
        }
    }

    post {
        always {
            sh 'rm -f config.yaml || true'
            cleanWs()
        }
        success {
            echo "✅ 构建成功！镜像: ${IMAGE_NAME}:${IMAGE_TAG}"
        }
        failure {
            echo "❌ 构建失败！"
        }
    }
}
