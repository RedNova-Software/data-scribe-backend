const AWS = require("aws-sdk");
const cognitoIdentityServiceProvider = new AWS.CognitoIdentityServiceProvider();

exports.handler = async (event, context, callback) => {
  console.log("Event: ", JSON.stringify(event, null, 2));

  if (event.triggerSource === "PostConfirmation_ConfirmSignUp") {
    const params = {
      UserPoolId: event.userPoolId,
      Username: event.userName,
    };

    try {
      await cognitoIdentityServiceProvider.adminDisableUser(params).promise();
      console.log("User disabled successfully");
    } catch (error) {
      console.error("Error disabling user: ", error);
      throw error;
    }
  }

  // Return to Amazon Cognito
  callback("Your account is pending admin approval.", event);
};