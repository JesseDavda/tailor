# Tailor - AI-Powered Resume Tailoring CLI

[![Tests](https://github.com/JesseDavda/tailor/actions/workflows/test.yml/badge.svg)](https://github.com/JesseDavda/tailor/actions/workflows/test.yml)

Tailor is a Go CLI tool that uses Claude API to intelligently adapt your master RenderCV YAML resume to specific job descriptions. It maintains complete factual accuracy while optimizing for ATS keywords and relevance.

This is **not** an automated application tool, tailor is built to give you a starting point to allow you to more efficiently show your expertise and personality to prospective employers.

AI can make mistakes, double check everything that is generated. All generated documents must go through a human reveiw process before applying for a position.

## Features

- **AI-Powered Tailoring**: Uses Claude API to reframe and re-prioritize your resume content
- **RenderCV Compatible**: Works seamlessly with RenderCV YAML format
- **Factual Accuracy**: Never invents or exaggerates - only reframes existing content
- **ATS Optimization**: Incorporates job description keywords naturally
- **Fast & Efficient**: Tailors resumes in seconds
- **Flexible Input**: Supports both job description files and direct text input
- **Dry Run Mode**: Preview changes before writing output
- **Verbose Mode**: Detailed progress reporting

## Installation

### Prerequisites

- Go 1.25.6 or higher
- Anthropic API key ([get one here](https://console.anthropic.com/))

### Build from Source

```bash
git clone <repository-url>
cd tailor
go build -o tailor cmd/tailor/main.go
```

### Set Up API Key

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```

Or pass it directly via the `--api-key` flag.

## Usage

### Basic Usage

Tailor your resume using a job description file:

```bash
./tailor --resume master.yaml --job job-desc.txt
```

### Direct Text Input

Use job text directly without creating a file:

```bash
./tailor --resume master.yaml --job-text "Senior Software Engineer at ACME Corp.
Requirements: 5+ years Go experience, distributed systems..."
```

### Custom Output Path

Specify where to save the tailored resume:

```bash
./tailor -r master.yaml -j job.txt -o acme-corp-resume.yaml
```

### Dry Run

Preview without writing output:

```bash
./tailor -r master.yaml -j job.txt --dry-run
```

### Verbose Mode

Get detailed progress information:

```bash
./tailor -r master.yaml -j job.txt -v
```

## Command Line Options

| Flag | Short | Description | Required |
|------|-------|-------------|----------|
| --resume | -r | Path to master RenderCV YAML resume | Yes |
| --job | -j | Path to job description file | * |
| --job-text | | Direct job description text input | * |
| --output | -o | Output path for tailored YAML | No (default: `tailored-resume.yaml`) |
| --cover-letter | -c | Generate cover letter with key qualifications and full text | No (default false) |
| --cover-letter-output | | Custom cover letter output path  | No (default: derived from resume output) |
| --api-key | | Anthropic API key | No (can use env var) |
| --model | -m | Anthropic model | No, will default to claude-haiku-4-5, the alias for the latest Haiku model |
| --verbose | -v | Enable verbose output | No |
| --dry-run | | Preview without writing output | No |
| --help | -h | Show help message | No |

\* Either `--job` or `--job-text` is required (not both)

## Environment Variables

- `ANTHROPIC_API_KEY`: Your Anthropic API key (alternative to `--api-key` flag)

## Example Workflow

1. **Prepare your master resume** in RenderCV YAML format
2. **Find a job posting** and save the description to a text file
3. **Tailor your resume**:
   ```bash
   tailor --resume ~/cv/master.yaml --job job-posting.txt
   ```
4. **Review the output**:
   ```bash
   cat tailored-resume.yaml
   ```
5. **Generate PDF with RenderCV**:
   ```bash
   rendercv render tailored-resume.yaml
   ```

## How It Works

1. **Reads** your master RenderCV YAML resume
2. **Analyzes** the job description for requirements and keywords
3. **Calls** Claude API with specialized prompts for resume tailoring
4. **Outputs** a tailored YAML maintaining the same structure
5. **Validates** the output is valid YAML before writing

### Tailoring Strategy

The tool uses carefully crafted prompts that instruct Claude to:

- Highlight relevant skills and experiences matching job requirements
- Reorder or emphasize achievements aligning with the role
- Incorporate keywords from the job description naturally
- Maintain chronological accuracy and factual correctness
- **Never** invent or exaggerate facts, skills, or experiences
- Keep all dates, company names, and titles accurate

## Limitations

- Requires valid RenderCV YAML format as input
- Requires internet connection for Claude API
- API calls cost tokens (typically $0.01-0.05 per resume)
- Output quality depends on master resume quality
- Does not generate PDFs (use RenderCV for that)

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes with clear commit messages
4. Add tests for new functionality
5. Submit a pull request

## License

The MIT License (MIT)

Copyright (c) 2026 Jesse Davda

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Powered by [Claude API](https://www.anthropic.com/claude)
- Compatible with [RenderCV](https://github.com/sinaatalay/rendercv)

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing issues for solutions
- Provide example inputs (redact personal information)

## Roadmap

Future enhancements:

- [ ] Unit and integration tests
- [ ] Caching for repeated tailoring
- [ ] Configuration file support (~/.tailorrc)
- [ ] Resume version tracking
- [ ] Diff view (master vs tailored)
- [ ] Cover letter generation
