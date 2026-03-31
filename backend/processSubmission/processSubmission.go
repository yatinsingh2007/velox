package processSubmission

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rishik92/velox/judge"
	"github.com/rishik92/velox/runBatch"
)

func ProcessSubmission(req judge.SubmissionRequest) judge.SubmissionResponse {
	var execCmd string
	var execArgs []string
	
	var filesToClean []string

	// 1. ROUTING: Compiled vs Interpreted
	switch req.Language {
	case "csharp":
		dirPath, dllPath, err := CompileInMemoryCSharp(req.SubmissionID, req.SourceCode)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = "dotnet"
		execArgs = []string{dllPath}
		filesToClean = append(filesToClean, dirPath)

	case "c":
		srcPath, binPath, err := CompileInMemoryC(req.SubmissionID, req.SourceCode)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = binPath
		execArgs = []string{}
		filesToClean = append(filesToClean, srcPath, binPath)

	case "cpp":
		srcPath, binPath, err := CompileInMemoryCPP(req.SubmissionID, req.SourceCode)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = binPath
		execArgs = []string{}
		filesToClean = append(filesToClean, srcPath, binPath)

	case "java":
		// Java requires defining a class name, usually Main. We create a directory for this submission.
		dirPath, className, err := CompileInMemoryJava(req.SubmissionID, req.SourceCode)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = "java"
		execArgs = []string{"-cp", dirPath, className}
		filesToClean = append(filesToClean, dirPath) // Clean up the entire submission directory

	case "python":
		scriptPath := fmt.Sprintf("/dev/shm/solution_%s.py", req.SubmissionID)
		if err := os.WriteFile(scriptPath, []byte(req.SourceCode), 0644); err != nil {
			return judge.SubmissionResponse{OverallState: "System Error: Cannot write to RAM"}
		}
		execCmd = "python3"
		execArgs = []string{scriptPath}
		filesToClean = append(filesToClean, scriptPath)

	case "node":
		scriptPath := fmt.Sprintf("/dev/shm/solution_%s.js", req.SubmissionID)
		if err := os.WriteFile(scriptPath, []byte(req.SourceCode), 0644); err != nil {
			return judge.SubmissionResponse{OverallState: "System Error: Cannot write to RAM"}
		}
		execCmd = "node"
		execArgs = []string{scriptPath}
		filesToClean = append(filesToClean, scriptPath)

	case "ts":
		jsPath, tsPath, err := CompileInMemoryTS(req.SubmissionID, req.SourceCode)
		if err != nil {
			return judge.SubmissionResponse{SubmissionID: req.SubmissionID, OverallState: "Compile Error", CompileError: err.Error()}
		}
		execCmd = "node"
		execArgs = []string{jsPath}
		filesToClean = append(filesToClean, tsPath, jsPath)

	default:
		return judge.SubmissionResponse{OverallState: "Unsupported Language"}
	}

	// 2. EXECUTION: Run the batch with the prepared command
	fmt.Printf("DEBUG: Executing Language: %s, Cmd: %s, Args: %v\n", req.Language, execCmd, execArgs)
	timeLimit := req.TimeLimitMs
	if timeLimit <= 0 {
		timeLimit = 3000 // 3 seconds default
	}
	memLimit := req.MemoryLimitKb
	if memLimit <= 0 {
		memLimit = 256000 // 256MB default
	}

	results := runBatch.RunBatch(execCmd, execArgs, req.TestCases, timeLimit, memLimit)

	// 3. CLEANUP: Delete files from /dev/shm RAM-disk
	for _, file := range filesToClean {
		os.RemoveAll(file)
	}

	// 4. AGGREGATE RESULTS
	overallState := "Accepted"
	for _, res := range results {
		if res.Status != "Accepted" {
			overallState = res.Status // e.g., "Wrong Answer" or "Time Limit Exceeded"
			// break // Fail-Fast: We found an error, no need to check the rest
		}
	}

	return judge.SubmissionResponse{
		SubmissionID: req.SubmissionID,
		OverallState: overallState,
		Results:      results,
	}
}

func CompileInMemoryC(submissionID, sourceCode string) (string, string, error) {
	// Use /tmp because /dev/shm is usually mounted as 'noexec' in Docker containers
	sourcePath := fmt.Sprintf("/tmp/solution_%s.c", submissionID)
	binaryPath := fmt.Sprintf("/tmp/solution_%s_c", submissionID)
	os.WriteFile(sourcePath, []byte(sourceCode), 0644)

	cmd := exec.Command("gcc", sourcePath, "-O2", "-o", binaryPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("compile error: %v, %s", err, string(out))
	}
	return sourcePath, binaryPath, nil
}

func CompileInMemoryCPP(submissionID, sourceCode string) (string, string, error) {
	sourcePath := fmt.Sprintf("/tmp/solution_%s.cpp", submissionID)
	binaryPath := fmt.Sprintf("/tmp/solution_%s_cpp", submissionID)
	os.WriteFile(sourcePath, []byte(sourceCode), 0644)

	cmd := exec.Command("g++", sourcePath, "-O2", "-o", binaryPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("compile error: %v, %s", err, string(out))
	}
	return sourcePath, binaryPath, nil
}

func CompileInMemoryJava(submissionID, sourceCode string) (string, string, error) {
	// Java requires the file name to match the public class name. We assume "Main".
	dirPath := fmt.Sprintf("/dev/shm/sol_%s", submissionID)
	os.MkdirAll(dirPath, 0755)
	sourcePath := fmt.Sprintf("%s/Main.java", dirPath)
	os.WriteFile(sourcePath, []byte(sourceCode), 0644)

	cmd := exec.Command("javac", sourcePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("compile error: %v, %s", err, string(out))
	}
	return dirPath, "Main", nil
}

func CompileInMemoryTS(submissionID, sourceCode string) (string, string, error) {
	sourcePath := fmt.Sprintf("/dev/shm/solution_%s.ts", submissionID)
	jsPath := fmt.Sprintf("/dev/shm/solution_%s.js", submissionID)
	os.WriteFile(sourcePath, []byte(sourceCode), 0644)

	cmd := exec.Command("npx", "tsc", sourcePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		// Output includes TS compile errors. It still might write a .js file, but we should fail the flow.
		return "", "", fmt.Errorf("compile error: %v, %s", err, string(out))
	}
	return jsPath, sourcePath, nil
}

func CompileInMemoryCSharp(submissionID, sourceCode string) (string, string, error) {
	dirPath := fmt.Sprintf("/tmp/sol_cs_%s", submissionID)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create temp dir: %v", err)
	}

	csprojContent := `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net8.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
  </PropertyGroup>
</Project>`
	if err := os.WriteFile(dirPath+"/project.csproj", []byte(csprojContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write csproj: %v", err)
	}

	if err := os.WriteFile(dirPath+"/Program.cs", []byte(sourceCode), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write source code: %v", err)
	}

	cmd := exec.Command("dotnet", "build", dirPath, "-c", "Release", "-o", dirPath+"/out", "-v", "q", "--nologo")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("compile error: %v, %s", err, string(out))
	}

	return dirPath, dirPath + "/out/project.dll", nil
}