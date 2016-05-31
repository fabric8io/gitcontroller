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
  def tag = 'latest'

  stage 'Stage'
  def stagedProject = stage()
  kubernetes.image().withName(imageName).build().fromPath(".")
  kubernetes.image().withName(imageName).tag().inRepository('docker.io/fabric8/'+imageName).force().withTag(tag)
  kubernetes.image().withName('docker.io/fabric8/'+imageName).push().withTag(tag).toRegistry()

  stage 'Promote'
  release(stagedProject)

}
def externalImages(){
  return ['git-controller']
}

def stage(){
  return stageProject{
    project = 'fabric8io/gitcontroller'
    useGitTagForNextVersion = true
    extraImagesToStage = externalImages()
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
