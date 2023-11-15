# GoBlog

GoBlog 是一个使用 Go 语言编写的简单博客项目，用于练习和学习 Go 语言。

## 安装

确保你的计算机上已经安装了 Go。然后，你可以通过以下命令来获取这个项目：

```bash
go get -u github.com/LuStormstout/goblog
```

## 使用

在项目的根目录下运行以下命令来启动服务器：

```bash
go run main.go
```

如果你在 macOS 上，并且希望使用 [air](https://github.com/cosmtrek/air) 来自动重载你的应用，你需要首先安装 air。你可以通过以下命令来安装 air：

```bash
brew install cosmtrek/air/air
```

安装完成后，你可以在项目的根目录下运行以下命令来启动服务器：

```bash
air
```

air 会监视你的代码文件，并在文件发生变化时自动重建并重启你的应用。

然后，你可以在浏览器中访问 `http://localhost:3000` 来查看你的博客。

## 贡献

欢迎任何形式的贡献。如果你发现了任何问题或者有任何建议，欢迎提交 issue 或者 pull request。

## 许可证

这个项目遵循 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。