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
