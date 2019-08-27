# oci-fn-bootstrapper
Function to bootstrap OCI setup for running FaaS functions. The bootstrapper function is idempotent and creates OCI resources only if they do not already exist.

# Prerequisites
* OCI Admin Credentials (Will be changed to resource principal in future)
  * User-OCID
  * Tenant-OCID
  * User-Fingerprint
  * User-Private-Key Base64 Encoded
  * Compartment-ID where the functions will be created
  * VCN Display Name
  * Regional Subnet Display Name
  * OCIR Repository Name
* fn-cli is installed (https://github.com/fnproject/cli)
* docker-cli is installed on your machine
  
# What does this function do?
This function, when invoked, ensures VCN/Subnet/Internet Gateway/Identity Policies/Security List/Repository are created. The output of the invocation contains the subnet/vcn OCID. The output also contains fn-cli context file contents in json format (minus user private key).


# How to use this?
* Assumption: The function code in this repository is already deployed to OCI FaaS under the application `bootstrap` and as function `oci-fn-bootstrap`.
* Assumption: The fn cli context bootstrapper exists (this context should belong to Admin who can create the fn setup)
* Prepare the config
```
{
  "user_id": "user-ocid",
  "tenant_id": "tenant-ocid",
  "region": "us-ashburn-1",
  "fingerprint": "user-key-fingerprint",
  "private_key": "Base64 Encoded RSA Private Key",
  "tenant_name": "tenant-name",
  "compartment_id": "compartment-ocid where the setup needs to be created",
  "vcn_name": "bootstrap-demo",
  "regional_subnet": "bootstrap-demo-regional",
  "repo_name": "bootstrap-demo"
}

```
# Invoke the function to create your setup 
    `DEBUG=1 cat config.json|fn --verbose --context bootstrapper invoke --display-call-id  bootstrapper oci-fn-bootstrapper`
    
# What do I do with the output?
Copy the fdk_context portion of the output into `~/.fn/contexts/bootstrap-demo.yaml` and format the json to yaml format. Point private-key in the file to your private-key, replace fingerprint/user-id if needed.

Now you can create an app using the following command
`fn --context bootstrap-demo create app --annotation oracle.com/oci/subnetIds='["subnet-id-from-the-earlier-output"]' demo1 `

Now you can create a function under the above app
`fn --context bootstrap-demo create function demo1 demo1 iad.ocir.io/odx-jafar/bootstrap-demo/some-image-you-have-pushed-to-the-repository`

Finally, invoke the function
`fn --context bootstrap-demo invoke demo1 demo1`

Successful invocation indicates that the setup was created successfully.