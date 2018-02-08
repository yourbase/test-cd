// Secrets for the GitHub CI. These are bogus strings. Before doing a production deployment,
// replace them with real ones.
//
// One pattern might be to change the path to the secrets in WORKSPACE from
// "secrets" to "secrets/prod" and use a script or git submodule to update that
// directory with real secrets when the deployment is being run from the right
// account.
//
// TODO: find a good pattern for replacing secrets based on the runtime context
// (dev, prod, etc).
{
	"username": std.base64("nictuku"),
	"password" : std.base64("bogus"),
	"token" : std.base64("bogus")
}
