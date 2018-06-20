# New Project Template

This repository contains a template you can use to seed a repository for a
new open source project.

See [releasing a project](https://gitlab.sas.com/techoffice/open-source-guide/blob/master/docs/creating/RELEASING.md) for more information about
releasing a new SAS open source project.

This template uses the [Apache 2.0 license](https://www.apache.org/licenses/LICENSE-2.0), which is the SAS default.  See the
documentation for instructions on using alternate license.

## How to use this template

1. Check this project out from GitLab.
    * There is no reason to fork it.
1. Create a new local repository and copy the files from this repo into it.
1. Develop your new project!

``` shell
git clone https://gitlab.sas.com/techoffice/new-project
mkdir my-new-thing
cd my-new-thing
git init
cp ../new-project/* .
git add *
git commit -a -m 'Boilerplate for new SAS open source project'
```

## Source Code Headers

Every file containing source code must include copyright and license
information.

### Example Apache license copyright header
Example copyright header for projects using the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) license:

    Copyright 2018 SAS Institute Inc.

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

### Example Creative Commons license copyright header
Example copyright header for projects using the [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) license: 

    Copyright 2018 SAS Institute Inc.

    This work is licensed under a Creative Commons Attribution 4.0
    International License. You may obtain a copy of the License at
    https://creativecommons.org/licenses/by/4.0/.