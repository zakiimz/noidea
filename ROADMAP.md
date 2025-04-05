## 🧠 `noidea` — Development Roadmap 2.0
**Vision:**
To create the most delightful, intelligent Git companion that makes every commit experience both productive and enjoyable.

#### Tasks:

Future Vision Idea:
- Integrate with GitHub Issues
- AI Powered GitHub Issue Creation
- Link commit messages to Issues and close them upon completion.
- Complete Issue tracking logic for your project within the terminal
    - Super easy to use with little commands.
    - Make it easy to check issues, work on them, commit changes, close issue. wokflow.
-> Expand current hook setup to seemlesly integrate these features

-> We could possibly link this with GitHub Projects. This way we can track Tasks we need to do, currently doing and done etc.
-> Need to explore how we could integrate with this nicely code whise


Another Idea:
-> Use AI to improve release notes.
    -> maybe there is a way where after automatic generation of relesae notes and "Whats Changed" listing all the commits.
    -> We could pass all that info to the LLM and let it generate a short overview section explaining what changed.

        This would make it a lot easier for non techincal users to understand what a version upgrade has in store for them.

        -> Only issue is how to hook into that process after automatic Changelog generation

## GitHub Integration and Authentication

### Implementation Plan

1. **GitHub Authentication**:
   - Extend the current secure storage system to handle GitHub Personal Access Tokens (PATs)
   - Create simple command to authenticate: `noidea github auth`
   - Store tokens securely using existing `secure` package
   - Validate tokens against GitHub API for immediate feedback

2. **Core Components**:
   - `internal/secure/github.go`: GitHub-specific validation and token handling
   - `internal/github/client.go`: GitHub API client implementation
   - `cmd/github.go`: CLI commands for GitHub integration

3. **User Experience**:
   - Simple one-time authentication flow: `noidea github auth`
   - Token stored securely in system keyring
   - Automatic token refresh when needed
   - Clear feedback on authentication status

4. **Implementation Benefits**:
   - Leverages existing secure storage system
   - Provides foundation for all GitHub-related features
   - Enables Post-Tag Hook and Release Notes enhancement
   - Future-proofs for GitHub Issues integration

5. **Security Considerations**:
   - Only request minimum required scopes for tokens
   - Support for token rotation and expiration
   - Clear documentation on token usage and security

## Enhanced Release Notes Feature

### Implementation Plan

1. **Post-Tag Hook Integration**:
   - Create a new `post-tag` Git hook that triggers when a version tag is created
   - Hook will gather commit information between the current tag and previous tag
   - Send collected commit data to our LLM to generate user-friendly release summary
   - Update GitHub release description via API with enhanced notes

2. **Core Components**:
   - `internal/hooks/posttag.go`: Hook handler for post-tag events
   - `internal/releaseai/`: New package to handle LLM interaction for release summaries
   - `internal/github/releases.go`: GitHub API integration for updating releases

3. **User Experience**:
   - User runs the existing `version.sh` script to bump version and create a tag
   - The post-tag hook automatically activates
   - Hook collects commit history, generates enhanced notes, and updates GitHub release
   - No additional steps required from the user

4. **Benefits**:
   - More control over the process than GitHub Actions approach
   - Runs on the developer's machine where context is available
   - Can be used regardless of CI/CD platform
   - Cleaner integration with existing version management workflow

5. **Implementation Phases**:
   - **Phase 1**: Basic post-tag hook that gathers commit information
   - **Phase 2**: LLM integration for generating human-readable summaries
   - **Phase 3**: GitHub release API integration to update release descriptions
   - **Phase 4**: Advanced customization options and templates

# NoIdea → AI Project Assistant: Integration Plan

## GitHub API Integration

- **GitHub REST/GraphQL API**: Implement a new `github` package in the `internal` directory to handle all GitHub API interactions.
  - Use GitHub's GraphQL API for efficient querying of issues, projects and complex relationships
  - Use REST API for simpler operations when appropriate

- **Authentication Options**:
  - GitHub Personal Access Tokens (PATs) stored securely using existing `secure` package
  - GitHub CLI authentication compatibility (gh auth)
  - OAuth flow for better security and scoped permissions

## Command Structure

Extend the current CLI with new commands:

```
noidea issue list              # List issues for current repo and their status + project relation
noidea issue view <number>     # View details of a specific issue
noidea issue create            # Create a new issue (with AI assistance)

noidea project init             # Initialize a new Project inside the current git repository. (Offer AI Assitance)
```

think about how to integrate into dev workflow so starting to work on issues and closing them with a commit flows and feels nice

## Core Components

1. **GitHub Module** (`internal/github/`):
   - Authentication handlers
   - Issue CRUD operations
   - Project boards interaction
   - Repository metadata management
   - Tracking relationships between commits and issues

2. **Workflow Engine** (`internal/workflow/`):
   - Branch naming conventions based on issues
   - Automated status transitions
   - Smart linking between commits and issues

3. **Enhanced AI Context** (`internal/context/`):
   - Expand existing AI capabilities to analyze issue descriptions
   - Generate issue summaries from code changes
   - Suggest issue status changes based on commit content

## Technical Implementation Details

1. **Local Cache and Syncing**:
   - Implement a local cache of GitHub issues/projects to improve performance
   - Use webhooks or polling to keep local state in sync
   - Store in SQLite or similar lightweight DB in user profile

2. **State Management**:
   - Track which issue is currently being worked on
   - Associate branches with issues
   - Maintain project board state locally

3. **Git Hooks Enhancement**:
   - Update existing hooks to detect issue references
   - Automatically add issue links to commits
   - Prompt for issue closure when appropriate

4. **AI Integration Points**:
   - Issue summarization (convert verbose issues to concise tasks)
   - Issue creation (generate well-formed issues from natural language)
   - Commit-to-issue matching (suggest which issue a commit addresses)
   - Work estimation (suggest timeframes based on issue complexity)

## User Experience Workflow

1. User runs `noidea issue list` to see pending tasks
2. User selects an issue with `noidea issue start #123`
   - Creates branch automatically
   - Updates issue status in GitHub Projects
3. User makes changes, commits with regular git workflow
   - The enhanced commit hook links to the issue
4. User completes work with `noidea issue close #123`
   - Tool generates AI summary of changes
   - Updates GitHub issue with relevant details
   - Moves card in project board

## Integration Plan Phases

1. **Phase 1**: Basic GitHub Issues integration
   - Authentication
   - Issue listing/viewing
   - Simple issue creation

2. **Phase 2**: Workflow integration
   - Branch management
   - Issue status transitions
   - Commit linking

3. **Phase 3**: Projects integration
   - Project board visualization
   - Card movement
   - Advanced AI assistance

4. **Phase 4**: Full workflow automation
   - Smart suggestions
   - Predictive issue management
   - Team collaboration features

This roadmap is a living document that will evolve based on community feedback and emerging priorities. The team welcomes suggestions and contributions to help shape the future of `noidea`.
