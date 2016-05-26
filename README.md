![Monkey-Ops logo](resources/images/monkey-ops-logo.jpg)

***

## What is Monkey-Ops

Monkey-Ops is a simple service implemented in Go which is introduced into a OpenShift V3.X and generates some chaos within it. Monkey-Ops seeks some Openshift components like Pods or DeploymentConfigs and randomly terminates them.


## Why Monkey-Ops

When you are implemented Cloud aware applications, these applications need to be designed so that they can tolerate the failure of services. Failures happen, and they inevitably happen when least desired, so the best way to prepare your application to fail is to test it in a chaos environment, and this is the target of Monkey-Ops.

Monkey-Ops is built to test the Openshift application's resilience, not to test the Openshift V3.X resilience.

## How to use Monkey-Ops

Monkey-Ops is prepared to running into a docker image. Monkey-Ops also includes an Openshift template in order to be deployed into a Openshift Project.

Monkey-Ops has two different modes of execution: background or rest.

* **Background**: With the Background mode, the service is running nonstop until you stop the container.
* **Rest**: With the Rest mode, you consume an api rest that allows you login in Openshift, choose a project, and execute the chaos for a certain time.

The service accept parameters as flags or environment variables. These are the input flags required:

      --API_SERVER string     API Server URL
      --INTERVAL float        Time interval between each actuation of operator monkey. It must be in seconds (by default 30)
      --MODE string           Execution mode: background or rest (by default "background")
      --PROJECT_NAME string   Project to get crazy
      --TOKEN string          Bearer token with edit grants to access to the Openshift project


