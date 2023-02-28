# pds-automation-demo

Note: Developed on ansible core 2.14.1 and python version = 3.11.1. Earlier versions may not function as expected.

Ansible playbooks to automate the deployment of PDS databases and px-delivery application.

Create an API key https://pds.docs.portworx.com/user-guide/add-api-keys/
Set it as the ```token``` variable in spin_pxdelivery.yaml and deploy_pxdelivery.yaml playbooks

Define the following variables
-  deployment_cluster:
-  pds_namespace:
-  my_ds_name:
-  pds_account:
-  pds_tenant:
-  pds_project:


The ansible fact `deploy` controls whether the PDS services are created and the application is deployed.
-  If set to `false` the playbook runs all the queries and builds the body_request for the POST command and writes it to disk.
-  If set to `true` it also runs the POST command and deploys the app

Run the playbook
```
ansible-playbook spin_pxdelivery.yaml
```

PDS API refernece https://prod.pds.portworx.com/swagger/index.html#/
