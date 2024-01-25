import path = require("path");
import fs = require("fs");
import { exec } from "child_process";

const lambdasSourceDir = path.join(__dirname, "../../../api/lambdas");
const outputDir = path.join(__dirname, "../../bin/lambdas");

// Ensure the output root directory exists
if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir, { recursive: true });
}

// Function to recursively find handler.go files
function findGoFiles(dir: string, fileList: string[] = []): string[] {
  const files: string[] = fs.readdirSync(dir);

  files.forEach((file) => {
    const filePath: string = path.join(dir, file);
    if (fs.lstatSync(filePath).isDirectory()) {
      findGoFiles(filePath, fileList);
    } else if (file === "handler.go") {
      fileList.push(filePath);
    }
  });

  return fileList;
}

const goFiles: string[] = findGoFiles(lambdasSourceDir);

goFiles.forEach((goFilePath) => {
  const folderPath: string = path.dirname(goFilePath);
  const parentFolder: string = path.basename(folderPath);
  const outputPath: string = path.join(outputDir, parentFolder);

  // Ensure the specific output directory exists
  if (!fs.existsSync(outputPath)) {
    fs.mkdirSync(outputPath, { recursive: true });
  }

  // Compile the Go file in an asynchronous way
  exec(
    `GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o ${path.join(
      outputPath,
      "bootstrap"
    )} ${goFilePath}`,
    {
      cwd: folderPath,
    },
    (error, stdout, stderr) => {
      if (error) {
        console.error(
          `Error compiling Go Lambda function in folder ${parentFolder}:`,
          error.message
        );
        return;
      }
      if (stderr) {
        console.error(`stderr: ${stderr}`);
        return;
      }
      console.log(`Successfully compiled Lambda function: ${parentFolder}`);
    }
  );
});
