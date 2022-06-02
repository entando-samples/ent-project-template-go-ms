# Deployment and installation
With this configuration, you can use the ent cli (https://dev.entando.org/next/docs/reference/entando-cli.html) to perform the full deployment sequence:

# Development notes
## Tips
* Check for the "CHANGE-IT" placeholders, replace it with your chosen docker image key

### Setup the project directory.
1. Prepare the bundle directory: `cp -r bundle_src bundle`
2. Initialize the project: `ent prj init`
3. Initialize publication: `ent prj pbs-init` (requires the git bundle repo url)
4. Attach to kubernetes for an Entando application via ent attach-kubeconfig config-file or similar

### Publish the bundle.
1. Build FE: ent prj fe-build -a
2. Build and publish BE: ./prepareDockerImage.sh
3. Publish FE: `ent prj fe-push`
4. Deploy (after connecting to k8s above): `ent prj deploy`
5. Install the bundle via 1) App Builder, 2) `ent prj install`, or 3) `ent prj install --conflict-strategy=OVERRIDE` on subsequent installs.
6. Iterate steps 1-4 to publish new versions.

## Local testing of the project
You can use the following commands from the application folder to run the local stack
* `ent prj xk start` - or stop to shutdown keycloak again.
* `cd src/main/go/` and `go run .` - to run the microservice
* `ent prj fe-test-run` - to run the React frontend

## Local setup
* Access Keycloak at http://localhost:9080/auth/
* Access swaggerui at http://localhost:8081/swagger-ui/
* Use client web_app when authorizing the microservices

## Notes
* Three users are included in the test keycloak realm config
  * user1/user1 - no roles
  * user2/user2 - first-role
  * user3/user3 - first-role second-role 
