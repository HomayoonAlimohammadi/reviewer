# AI Code Reviewer

## Overview

AI Code Reviewer is a tool designed to assist developers in reviewing code by leveraging artificial intelligence. It helps identify potential issues, suggests improvements, and ensures code quality and consistency.

## Features

- Automated code review using AI
- Supports multiple programming languages
- Provides suggestions for code improvements
- Integrates with popular version control systems
- Customizable rules and configurations

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [Docker](https://docs.docker.com/get-docker/) (optional, for running in a containerized environment)
- [`make`](https://www.gnu.org/software/make/)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/HomayoonAlimohammadi/reviewer.git
    cd reviewer
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

### Usage

To start the AI Code Reviewer, run the following command:

```sh
make run
```

To cleanup everything:

```sh
make clean
```
