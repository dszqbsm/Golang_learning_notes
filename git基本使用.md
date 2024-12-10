# 将本地项目代码提交到远程仓库

1. 创建远程仓库

2. 初始化本地git仓库：在项目目录下执行`git init`

> 远程仓库名和本地仓库名不需要相同

3. 添加忽略文件：在项目目录下创建`.gitignore`文件，并添加忽略文件，比如忽略历史修改记录，在`.gitignore`文件中输入`.history/`

4. 暂存所有更改：`git add .`

5. 将暂存的更改提交到本地代码仓库并编写提交说明：`git commit -m "提交说明"`

6. 将本地git仓库和远程github仓库关联：`git remote add origin git@github.com:dszqbsm/test.git`

> 在远程仓库链接后加`.git`

7. 将本地代码仓库推送到远程仓库：`git push -u origin master`

> 默认分支名称为master，可以使用`git branch -M main`将master分支重命名为main



# 忽略某些文件

1. 要忽略的文件或目录没有被提交推送过

在项目初始化时便创建`.gitignore`文件，并在其中写好要忽略的文件或目录

2. 要忽略的文件或目录已经被提交推送过

> 因为`.gitignore`文件只能忽略那些原来没有被提交推送过的文件，如果某些文件已经被纳入了版本管理中，那么修改`.gitignore`文件是无效的，必须删除文件的追踪

如执行`git rm --cached test/test.txt`表示删除对test.txt文件的追踪，但不会删除文件的数据，或者`git rm -r --cached .history`表示递归的删除对`.history`目录下的所有文件和目录的版本管理追踪

> 似乎直接在`.gitignore`文件中添加`.history/`并提交推送即可取消对history文件夹的修改追踪
