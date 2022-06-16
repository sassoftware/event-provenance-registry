# New Project Template Instructions

This repository contains a template you can use to seed a repository for a
new open source project, following the [SAS Open Source Guidelines](https://gitlab.sas.com/techoffice/open-source-guide/blob/master/README.md). See [releasing a project](https://gitlab.sas.com/techoffice/open-source-guide/blob/master/docs/creating/RELEASING.md) for more information about
releasing a new SAS open source project.

This template uses the [Apache 2.0 license](https://www.apache.org/licenses/LICENSE-2.0), which is the SAS default.  See the
documentation for instructions on using alternate license.

## How to use this template

1. Navigate to https://gitlab.sas.com/projects/new
1. Select the **Import project** tab.
1. Choose to import from **Repo by URL**.
1. For the **Git repository URL** enter https://gitlab.sas.com/techoffice/new-project.git
1. For the **Project name** enter the desired name for your project.
1. By default, the project will be created in your user space. If you want to create your project in a different location, specify that as part of the **Project URL**.
1. (Optional) Supply a project description.
1. Specify the **Visibility Level**. If you don't choose _Public_ (meaning open to all logged in SAS users) then you will need to provide access to any reviewers.
1. When you have supplied all the information, click the **Create project** button to create the project.

At this point, you will have a new Git repository that is populated with all of the new project template content.

**NOTE:** You should delete the INSTRUCTIONS.md file from your project before committing your changes.

## Source Code Headers

Every file containing source code must include copyright and license information. Source code refers to any executable code such as .java or .go files, shell scripts, etc. It does not refer to content such as documentation, build scripts, or configuration files.

**NOTE:** The required source code headers must be placed at the top of the source code files, above any other header information. Per SAS legal, "We're seeing increased reliance both at SAS and across the industry on automated code scanning and license management. Scanning tools expect to see copyright notices at the top."

### Example Apache license copyright header
Example copyright header for projects using the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) license.

#### New short form (preferred)

    Copyright © 2020, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
    SPDX-License-Identifier: Apache-2.0


#### Longer form allowed for backward compatibility

    Copyright © 2020, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.