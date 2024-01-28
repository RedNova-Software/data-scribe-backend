import { execSync } from "child_process";
import * as path from "path";
import * as fs from "fs";

export function buildLambda(): void {
  // Define the source and destination paths
  const sourcePath = path.join(__dirname, "../lambda");
  const outputPath = path.join(sourcePath, "/build/lambda.zip");

  // Remove the existing zip file if it exists
  if (fs.existsSync(outputPath)) {
    fs.unlinkSync(outputPath);
  }

  // Navigate to the lambda directory and zip its contents
  const zipCommand = `cd ${sourcePath} && zip -r ${outputPath} .`;

  try {
    execSync(zipCommand);
    console.log("Post Confirmation Lambda zipped successfully.");
  } catch (error) {
    console.error(`Error zipping the post confirmation lambda: ${error}`);
    throw error; // Rethrow the error to handle it further up the call stack if necessary
  }
}
