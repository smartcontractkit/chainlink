aws  --profile  EKS-ADMIN sso login
Attempting to automatically open the SSO authorization page in your default browser.
If the browser does not open or you wish to use a different device to authorize this request, open the following URL:

https://device.sso.us-west-2.amazonaws.com/

Then enter the code:

ZZKT-KVJT
Successfully logged into Start URL: https://smartcontract.awsapps.com/start
aws  --profile  PowerUserAccess-795953128386 sso login
aws  --profile  PowerUserAccess-795953128386 ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 795953128386.dkr.ecr.us-west-2.amazonaws.com
