# OpenShift Template Processor

This project supports [OpenShift Templates](https://docs.openshift.com/container-platform/4.2/openshift_images/using-templates.html).

The spec for OpenShift Templates as a resource type looks like:

```
- name: <Resource name>
  openshiftTemplate:
    template: <path to file, directory, or URL> # Required
    params: # Optional
      <key>: <value>
      ...
    paramFiles: # Optional
    - <path to file or directory>
    - ...
    paramDir: <path to directory> # Optional
```

## Processing Modes

OpenShift Templates can be processed in one of four modes, defined by the various combinations of templates and parameters.

## Single Template, Single Parameter Set

When `.openshiftTemplate.template` points to a file and `.openshiftTemplate.paramDir` is empty, we will process that template, passing the parameters to it.

```
- name: <Resource name>
  openshiftTemplate:
    template: templates/app-stack.yaml
    params:
      APP_NAME: myapp
      IMAGE: myorg/app:v1.0
    paramFiles:
    - myparams/common-params
    - myparams/params-dev
```

Becomes:

```
oc process -f templates/app-stack.yaml -p APP_NAME=myapp -p IMAGE=myorg/app:v1.0 --param-file=myparams/common-params --param-file=myparams/params-dev
```

## Multiple Templates, Single Parameter Set

When `.openshiftTemplate.template` points to a directory and `.openshiftTemplate.paramDir` is empty, we will process all templates in that directory, passing the parameters to each template.

```
- name: <Resource name>
  openshiftTemplate:
    template: templates/
    params:
      APP_NAME: myapp
      IMAGE: myorg/app:v1.0
    paramFiles:
    - myparams/common-params
    - myparams/params-dev
```

Becomes:

```
oc process -f templates/app-stack1.yaml -p APP_NAME=myapp -p IMAGE=myorg/app:v1.0 --param-file=myparams/common-params --param-file=myparams/params-dev
oc process -f templates/app-stack2.yaml -p APP_NAME=myapp -p IMAGE=myorg/app:v1.0 --param-file=myparams/common-params --param-file=myparams/params-dev
oc process -f templates/app-stack3.yaml -p APP_NAME=myapp -p IMAGE=myorg/app:v1.0 --param-file=myparams/common-params --param-file=myparams/params-dev
...
```

## Single Template, Multiple Parameter Sets

When `.openshiftTemplate.template` points to a file and `.openshiftTemplate.paramDir` is points to a directory, we will process the template once for each parameter file in that directory. Additionally, the parameters in `.openshiftTemplate.params` will be passed to each template. `.openshiftTemplate.paramFiles` will be ignored.

```
- name: <Resource name>
  openshiftTemplate:
    template: templates/app-stack.yaml
    params:
      IMAGE: myorg/app:v1.0
    paramDir: myparams/
```

Becomes:

```
oc process -f templates/app-stack.yaml -p IMAGE=myorg/app:v1.0 --param-file=myparams/params1
oc process -f templates/app-stack.yaml -p IMAGE=myorg/app:v1.0 --param-file=myparams/params2
oc process -f templates/app-stack.yaml -p IMAGE=myorg/app:v1.0 --param-file=myparams/params3
...
```

## Multiple Templates, Multiple Parameter Sets

When `.openshiftTemplate.template` points to a directory and `.openshiftTemplate.paramDir` is points to a directory, we will process each template in the `.template` directory, expecting a parameter file of the same name (minus the file extension) in the `.paramDir` directory. Additionally, the parameters in `.openshiftTemplate.params` will be passed to each template. `.openshiftTemplate.paramFiles` will be ignored.

```
- name: <Resource name>
  openshiftTemplate:
    template: templates/
    params:
      IMAGE: myorg/app:v1.0
    paramDir: myparams/
```

Becomes:

```
oc process -f templates/app-stack1.yaml -p IMAGE=myorg/app:v1.0 --param-file=myparams/app-stack1
oc process -f templates/app-stack2.yaml -p IMAGE=myorg/app:v1.0 --param-file=myparams/app-stack2
oc process -f templates/app-stack3.yaml -p IMAGE=myorg/app:v1.0 --param-file=myparams/app-stack3
...
```
