import { execSync } from 'child_process';
import path = require('path');
import fs = require('fs');

const lambdasSourceDir = path.join(__dirname, '../../../api/lambdas');
const outputDir = path.join(__dirname, '../../bin/lambdas');


// Ensure the output directory exists
if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir, { recursive: true });
}

// Read the directories under the source directory
const lambdaFolders = fs.readdirSync(lambdasSourceDir);

// Iterate over each folder and build the Go binaries
lambdaFolders.forEach(folder => {
  const folderPath = path.join(lambdasSourceDir, folder);
  const stat = fs.lstatSync(folderPath);

  if (stat.isDirectory()) {
    try {
      const goFilePath = path.join(folderPath, 'handler.go');
      const outputPath = path.join(outputDir, folder);
      const zipPath = `${outputPath}.zip`;

      // Compile the Go file
      execSync(`GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o ${outputPath}/bootstrap ${goFilePath}`, {
        stdio: 'inherit',
        cwd: folderPath,
      });

      console.log(`Successfully compiled Lambda function: ${folder}`);
    } catch (error: any) {
      console.error(`Error compiling Go Lambda function in folder ${folder}:`, error);
    }
  }
});