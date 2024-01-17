import { exec, execSync } from "child_process";
import path = require("path");
import fs = require("fs");

const lambdasSourceDir = path.join(__dirname, "../../../api/lambdas");
const outputDir = path.join(__dirname, "../../bin/lambdas");

// Ensure the output directory exists
if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir, { recursive: true });
}

const lambdaFolders = fs.readdirSync(lambdasSourceDir).filter((folder) => {
  const folderPath = path.join(lambdasSourceDir, folder);
  return fs.lstatSync(folderPath).isDirectory();
});

lambdaFolders.forEach((folder) => {
  const folderPath = path.join(lambdasSourceDir, folder);
  const goFilePath = path.join(folderPath, "handler.go");
  const outputPath = path.join(outputDir, folder);

  // Compile the Go file in an asynchronous way
  exec(
    `GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o ${outputPath}/bootstrap ${goFilePath}`,
    {
      cwd: folderPath,
    },
    (error, stdout, stderr) => {
      if (error) {
        console.error(
          `Error compiling Go Lambda function in folder ${folder}:`,
          error
        );
        return;
      }
      if (stderr) {
        console.error(`stderr: ${stderr}`);
        return;
      }
      console.log(`Successfully compiled Lambda function: ${folder}`);
    }
  );
});
