module.exports = {
    extends: ['@commitlint/config-conventional'],
    parserPreset: 'conventional-changelog-conventionalcommits',
    rules: {
        'body-leading-blank': [2, 'always'],
        'body-case': [
            1,
            'always',
            [
                'lower-case', // lower case
                'upper-case', // UPPERCASE
                'camel-case', // camelCase
                'kebab-case', // kebab-case
                'pascal-case', // PascalCase
                'sentence-case', // Sentence case
                'snake-case', // snake_case
                'start-case', // Start Case
            ]
        ],
        'body-empty': [2, 'never'],
        'footer-leading-blank': [2, 'always'],
        'type-case': [
            1,
            'always',
            [
                'lower-case', // lower case
                'upper-case', // UPPERCASE
                'camel-case', // camelCase
                'kebab-case', // kebab-case
                'pascal-case', // PascalCase
                'sentence-case', // Sentence case
                'snake-case', // snake_case
                'start-case', // Start Case
            ]
        ],
        'type-enum': [
            2, // Error if the type is not in the enum
            'always', // Type is always required
            [
                'init',
                'build',
                'chore',
                'ci',
                'docs',
                'feat',
                'fix',
                'perf',
                'refactor',
                'revert',
                'style',
                'test',
            ],
        ],
        'scope-case': [
            1,
            'always',
            [
                'lower-case', // lower case
                'upper-case', // UPPERCASE
                'camel-case', // camelCase
                'kebab-case', // kebab-case
                'pascal-case', // PascalCase
                'sentence-case', // Sentence case
                'snake-case', // snake_case
                'start-case', // Start Case
            ]
        ],
        'scope-empty': [2, 'never'],
        'scope-enum': [
            2,
            'always',
            [
                'packages',
                'core',
                'entities',
                'setup',
                'repo'
            ]
        ],
        'scope-max-length': [2, 'always', 30],
        'signed-off-by': [2, 'always'],
        'subject-empty': [2, 'never'],
        'subject-full-stop': [2, 'never'],
        'subject-max-length': [2, 'always', 72]
    },
    prompt: {
        questions: {
            type: {
                description: "Select the type of change that you're committing",
                enum: {
                    init: {
                        description: 'Changes that initialize a component',
                        title: 'Init',
                        emoji: '⏩'
                    },
                    build: {
                        description:
                            'Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)',
                        title: 'Builds',
                        emoji: '🛠',
                    },
                    chore: {
                        description: "Other changes that don't modify src or test files",
                        title: 'Chores',
                        emoji: '♻️',
                    },
                    ci: {
                        description:
                            'Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)',
                        title: 'Continuous Integrations',
                        emoji: '⚙️',
                    },
                    docs: {
                        description: 'Documentation only changes',
                        title: 'Documentation',
                        emoji: '📚',
                    },
                    feat: {
                        description: 'A new feature',
                        title: 'Features',
                        emoji: '✨',
                    },
                    fix: {
                        description: 'A bug fix',
                        title: 'Bug Fixes',
                        emoji: '🐛',
                    },
                    perf: {
                        description: 'A code change that improves performance',
                        title: 'Performance Improvements',
                        emoji: '🚀',
                    },
                    refactor: {
                        description:
                            'A code change that neither fixes a bug nor adds a feature',
                        title: 'Code Refactoring',
                        emoji: '📦',
                    },
                    revert: {
                        description: 'Reverts a previous commit',
                        title: 'Reverts',
                        emoji: '🗑',
                    },
                    style: {
                        description:
                            'Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)',
                        title: 'Styles',
                        emoji: '💎',
                    },
                    test: {
                        description: 'Adding missing tests or correcting existing tests',
                        title: 'Tests',
                        emoji: '🚨',
                    },
                },
            },
            scope: {
                description:
                    'What is the scope of this change (e.g. component or file name)',
            },
            subject: {
                description:
                    'Write a short, imperative tense description of the change',
            },
            body: {
                description: 'Provide a longer description of the change',
            }
        },
    },
};
