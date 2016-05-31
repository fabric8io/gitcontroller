#!/usr/bin/groovy

node{

  checkout scm

  kubernetes.pod('buildpod').withImage('fabric8/go-builder')
  .withEnvVar('GOPATH','/home/jenkins/workspace/workspace/go')
  .withPrivileged(true).inside {

    stage 'build binary'

    sh "mkdir -p ../go/src/github.com/fabric8io/gitcontroller; cp -R ../${env.JOB_NAME}/. ../go/src/github.com/fabric8io/gitcontroller/; cd ../go/src/github.com/fabric8io/gitcontroller; make"

    sh "cp -R ../go/src/github.com/fabric8io/gitcontroller/build ."
  }
  def imageName = 'gitcontroller'

  stage 'Stage'
  def stagedProject = stage()
  def tag = stagedProject[1]

  kubernetes.image().withName(imageName).build().fromPath(".")
  kubernetes.image().withName(imageName).tag().inRepository("${env.FABRIC8_DOCKER_REGISTRY_SERVICE_HOST}:${env.FABRIC8_DOCKER_REGISTRY_SERVICE_PORT}/fabric8/"+imageName).force().withTag(tag)
  kubernetes.image().withName("${env.FABRIC8_DOCKER_REGISTRY_SERVICE_HOST}:${env.FABRIC8_DOCKER_REGISTRY_SERVICE_PORT}/fabric8/"+imageName).push().withTag(tag).toRegistry()


  stage 'Promote'
  release(stagedProject)

}
def externalImages(){
  return ['gitcontroller']
}

def stage(){
  return stageProject{
    project = 'fabric8io/gitcontroller'
    useGitTagForNextVersion = true
  }
}

def release(project){
  releaseProject{
    stagedProject = project
    useGitTagForNextVersion = true
    helmPush = false
    groupId = 'io.fabric8'
    githubOrganisation = 'fabric8io'
    artifactIdToWatchInCentral = 'git-controller'
    artifactExtensionToWatchInCentral = 'jar'
    promoteToDockerRegistry = 'docker.io'
    dockerOrganisation = 'fabric8'
    extraImagesToTag = externalImages()
  }
}
