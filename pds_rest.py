import requests
import json
import os

baseURL = "https://prod.pds.portworx.com/api"
#baseURL = "https://staging.pds.portworx.com/api"
bearer_token = str(os.getenv("BEARER_TOKEN"))
#bearer_token = str(os.getenv("STAGING_TOKEN"))

def bearer_oauth(r):
    """
    Method required by bearer token authentication.
    """

    r.headers["Authorization"] = f"Bearer {bearer_token}"
    r.headers["accept"] = "application/json"
    return r

def get_accounts():
    response = requests.get(
        baseURL + "/accounts", auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get rules (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def get_account(acc_id):
    response = requests.post(
        baseURL + "/accounts/" + acc_id, auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot create account (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def create_account():
    print("not yet implemented, I don't have permissions")
    # response = requests.post(
    #     baseURL + "/accounts", auth=bearer_oauth
    # )
    # if response.status_code != 200:
    #     raise Exception(
    #         "Cannot get rules (HTTP {}): {}".format(response.status_code, response.text)
    #     )
    # return response.json()

def accept_eula(acc_id):
    response = requests.put(
        baseURL + "/accounts/" + acc_id + "/eula", auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot create account (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def get_users(acc_id):
    response = requests.get(
       baseURL + "/accounts/" + acc_id + "/users", auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def get_dataservices():
    response = requests.get(
       baseURL + "/data-services/", auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def get_tenants(acc_id):
    response = requests.get(
       baseURL + "/accounts/" + acc_id + "/tenants" , auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def get_tenant(ten_id):
    response = requests.get(
       baseURL + "/tenants/" + ten_id, auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def get_tenant_dns(ten_id):
    response = requests.get(
       baseURL + "/tenants/" + ten_id + "/dns-details", auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def get_appconfig_template(app_id):
    response = requests.get(
       baseURL + "/application-configuration-templates/" + app_id, auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def get_appconfig_templates_all(ten_id):
    response = requests.get(
       baseURL + "/tenants/" + ten_id + "/application-configuration-templates/", auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()


def get_tenant_projects(tenant_id):
    response = requests.get(
       baseURL + "/tenants/" + tenant_id + "/projects", auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def get_deploymenttargets(tenant_id):
    response = requests.get(
       baseURL + "/tenants/" + tenant_id + "/deployment-targets", auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def delete_deploymenttargets(target_id):
    response = requests.delete(
       baseURL + "/deployment-targets/" + target_id, auth=bearer_oauth
    )
    print(response.status_code)
    if response.status_code != 204:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    # return response.json()

def get_deployments(project_id):
    response = requests.get(
       baseURL + "/projects/" + project_id + "/deployments", auth=bearer_oauth
    )
    if response.status_code != 200:
        raise Exception(
            "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
        )
    return response.json()

def delete_deployments(deployment_id):
    response = requests.delete(
       baseURL + "/deployments/" + deployment_id + "?force=true", auth=bearer_oauth
    )
    print('response:', response)
    print('response.txt', response.text)
    if response.status_code != 202:
         raise Exception(
             "Cannot get users (HTTP {}): {}".format(response.status_code, response.text)
         )
    return response.text



def unhealthy_cluster(target_json):
    item_dict = json.loads(target_json)
    
    items = len(item_dict)
    del_list = []
    print('Target K8s Cluster Status in PDS')
    for x in range(items):
        x_status = item_dict[x]["status"]
        x_name = item_dict[x]["name"]
        print(x_name + "status is " + x_status)
    print("...")     
    for x in range(items):
        x_status = item_dict[x]["status"]
        x_name = item_dict[x]["name"]
        x_id = item_dict[x]["id"]
        if x_status == "unhealthy":
            print("Time to DELETE " + x_name + " " + x_id )
            del_list.append(x_id)
    return del_list

def unhealthy_deployments(ten_id, dep_target_id):
    ten_proj = get_tenant_projects(ten_id)
    x = json.dumps(ten_proj['data'][0]['id'])
    deployments = get_deployments(x)
    #print(json.dumps(deployments, indent=4))
    items = len(deployments['data'])
    dep_list = []
    for x in range(items):
         x_deployment_target_id = deployments['data'][x]["deployment_target_id"]
         #print(x_deployment_target_id)
         x_dep_id = deployments['data'][x]['id']
         #print(x_dep_id)
         if x_deployment_target_id == dep_target_id:
            #print("MATCH " + x_dep_id + "deployment to be deleted on target " + x_deployment_target_id)
            dep_list.append(x_dep_id)
    return dep_list

        


def clean_unhealthy_clusters(target_json, ten_id):
    f = unhealthy_cluster(target_json)
    for x in f:
        g = unhealthy_deployments(ten_id, x)
        for y in g:
            delete_deployments(y)
        delete_deploymenttargets(x)


