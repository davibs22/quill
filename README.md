<div align="center">
    <img alt="quill-logo" height="100px" src="./assets/quill-logo.png">
</div>

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Generate meaningful Git commit messages with AI effortlessly.

## Overview
Quill is an open-source CLI tool that uses artificial intelligence to generate clear, concise, and context-aware Git commit messages. Say goodbye to git commit -m "fix" and let AI do the heavy lifting!

**Quill provides:**
- A command-line interface for generating commit messages.
- Integration with the OpenAI API for AI-powered commit message generation.
- Support for local models via Ollama for private, offline usage.
- Easy-to-use interface for generating commit messages.

## Usage
1. Stage your changes:
```bash
git add .
```
2. Generate a commit message:
```bash
quill
```
3. Commit your changes:
```bash
git commit -m "your commit message"
```
