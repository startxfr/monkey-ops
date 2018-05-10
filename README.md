![Monkey-Ops logo](resources/images/logo.svg)


## What is Monkey-Ops

Monkey-Ops is a simple service implemented in Go which is deployed into a OpenShift v3.x and generates some chaos within it. Monkey-Ops seeks some Openshift components like Pods or DeploymentConfigs and randomly terminates them.

## Why Monkey-Ops

When you are implemented Cloud aware applications, these applications need to be designed so that they can tolerate the failure of services. Failures happen, and they inevitably happen when least desired, so the best way to prepare your application to fail is to test it in a chaos environment, and this is the target of Monkey-Ops.

Monkey-Ops is built to test the Openshift application's resilience, not to test the Openshift V3.X resilience.

## How to use Monkey-Ops

Monkey-Ops is prepared to running into a docker image. Monkey-Ops also includes an Openshift template in order to be deployed into a Openshift Project.

Monkey-Ops has two different modes of execution: background or rest.

* **Background**: With the Background mode, the service is running nonstop until you stop the container.
* **Rest**: With the Rest mode, you consume an api rest that allows you login in Openshift, choose a project, and execute the chaos for a certain time.

Monkey-Ops has tree differents chaos method: pod, dc or random.

* **pod**: This method tell the agent to use only pod deletion for perturbing the service or application
* **dc**: This method tell the agent to use only deployementConfig change (reducing replica number) for perturbing the service or application
* **random**: This method tell the agent to use randomly alternate between theses two method for perturbing the project components

The service accept parameters as flags or environment variables. These are the input flags required:

      --API_SERVER string     API Server URL (use KUBERNETES_SERVICE_HOST if not provided)
      --INTERVAL float        Time interval between each actuation of operator monkey. It must be in seconds (by default 30)
      --MODE string           Execution mode: background or rest (by default "background")
      --METHOD string         Chaos method: pod, dc or random (default "random")
      --PROJECT_NAME string   Project to get crazy
      --TOKEN string          Bearer token with edit grants to access to the Openshift project
      
### Usage with Docker

#### Downloading the image

```bash
docker pull startx/monkey-ops:latest
```

#### Running the image

```bash
docker run startx/monkey-ops \
       /monkey-ops \
       --TOKEN="my-bearer-token-from-openshift-authentication" \
       --PROJECT_NAME="my-project" \
       --API_SERVER="https://mycluster.openshift.example.com:8443" \
       --INTERVAL=30 \
       --MODE=background \
       --METHOD=pod
```

### Usage with Openshift v3.x

#### Prepare chaos environement

Before all is necessary to connect to your openshift cluster and enter your target project

```bash
# Connecting to the openshift cluster
# <user> is your openshift username
# <password> is your openshift password
# <url> is your openshift manager url (usually on port 8443)
# ex: oc login -u dev -p dev https://openshift.demo.startx.fr:8443
oc login -u <user> -p <password> <url>
# Selecting your project to challenge
oc project demo
```

Monkey-ops require to create a service account (and a token as a secret) with editing permissions within the project that you want to use. 
The service account must be called with the same name than [mokey-ops template](https://raw.githubusercontent.com/startxfr/monkey-ops/master/openshift/monkey-ops-template.yml) parameter SA_NAME, by default monkey-ops.
This service account must also have the edit role for this project.

```bash
# Create the service account "monkey-ops"
oc create -f https://raw.githubusercontent.com/startxfr/monkey-ops/master/openshift/monkey-ops-sa.yml
# Grant the edit role for this project service account
# oc policy add-role-to-user edit system:serviceaccount:<project>:<service-account>
oc policy add-role-to-user edit system:serviceaccount:demo:monkey-ops
```

This will create a service account named *monkey-ops* and grant it with edit role

If you can't create this service account due to restricted permissions, you can find more information on how to do it in [Managing Service Accounts link](https://docs.openshift.com/enterprise/3.1/dev_guide/service_accounts.html#managing-service-accounts)

#### Deploy the monkey agent

Deploy the monkey agent template in your project using the [mokey-ops template](https://raw.githubusercontent.com/startxfr/monkey-ops/master/openshift/monkey-ops-template.yml)

```bash
oc create \
   -f https://raw.githubusercontent.com/startxfr/monkey-ops/master/openshift/monkey-ops-template.yml \
   -n "demo"
```

Then use this template to start your monkey-ops agent in your project

```bash
oc process monkey-ops-template \
   -p APP_NAME=monkey-ops \
   -p INTERVAL=30 \
   -p MODE=background \
   -p METHOD=pod \
   -p TZ=Europe/Paris | \
oc create -f -
```
	
Once you have monkey-ops running in your project, you can see what the service is doing in your application logs

```bash
oc logs po/monkey-ops
```

### API REST

Monkey-Ops Api Rest expose two endpoints:

* **/login**

>This endpoint allows a user to log into Openshift in order to get a token and  projects to which it belongs.

	
>**Request Input JSON:**


>{
>     "user": "User name",
>     "password": "User password",
>     "url": "Openshift API Server URL. e.g. https://ose.api.server:8443"
> }

>**Request Output JSON:**

>	{
>     "token": "Token",
>     "projects": {
>    	 "project1 name",
>    	 "project2 name",
>    	 .
>    	 .
>    	 .
>    	 "projectN name"
>    	 }
>}	 

	
* **/chaos**

>This endpoint allows a user to launch the monkey-ops agent for a certain time.

>**Request Input JSON:**

>	{
>     "token": "Token",
>     "url": "Openshift API Server URL. e.g. https://ose.api.server:8443",
>     "project": "Project name",
>     "interval": Time interval between each actuation in seconds,
>     "totalTime": Total Time of monkey-ops execution in seconds
>	}

## Credit 

The main project of this monkey-ops implementation came from [Produban monkey-ops project](https://github.com/Produban/monkey-ops)
