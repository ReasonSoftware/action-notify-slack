## Description

Please include a summary of the change and which issue is fixed. Please also include relevant motivation and context. List any dependencies that are required for this change.

Fixes # (issue)

## Type of change

Please delete options that are not relevant.

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] This change requires a documentation update

## How Has This Been Tested

Please describe the tests that you ran to verify your changes. Provide instructions so we can reproduce. Please also list any relevant details for your test configuration

- [ ] Test A
- [ ] Test B

#### Example Workflow

<details><summary>Clich Here to Expand</summary>

```yaml
name: release
on:
  push:
    branches:
      - master

jobs:
  build:
    name: build
    runs-on: ubuntu-latest
    steps:
      - run: |
          mkdir .github && touch attachment.json
          var='{"color": "#000000","fields": [{"title": "Repository","value": "<https://github.com/ReasonSoftware/PROJECT|PROJECT>","short": true},{"title": "Workflow","value": "<https://github.com/ReasonSoftware/PROJECT/actions?query=workflow%3AWORKFLOW|WORKFLOW>","short": true},{"title": "Initiator","value": "<https://github.com/USER|USER>","short": true},{"title": "Status","value": "STATUS","short": true}],"footer": "<https://github.com/ReasonSoftware/action-notify-slack|ReasonSoftware/action-notify-slack>","footer_icon": "https://cdn.reasonsecurity.com/images/logo.png","ts": "TIMESTAMP"}'
          echo $var | 
            sed "s/000000/fbf000/" | 
            sed "s/PROJECT/$(echo $GITHUB_REPOSITORY | awk -F'/' '{print $2}')/g" |
            sed "s/WORKFLOW/release/g" |
            sed "s/USER/$GITHUB_ACTOR/g" |
            sed "s/STATUS/STARTED/g" |
            sed "s/TIMESTAMP/$(date +%s)/g" > attachment.json

      - name: Notification
        uses: docker://reasonsoftware/action-notify-slack:latest
        env:
          TOKEN: ${{ secrets.SLACK_TOKEN }}
          CHANNEL: ${{ secrets.SLACK_CHANNEL }}
          ATTACHMENTS_FILE: .github/attachments.json
```

</details>

# Checklist

- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
