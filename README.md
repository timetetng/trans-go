# ✏️trans-go

一个基于 LLM 的命令行翻译工具，支持文本翻译、编程变量命名及单词查询，让你无需离开终端就能享受翻译，提高 coding 连贯性。

## 功能特性

* **文本翻译**：快速将文本翻译为目标语言（默认中译英）。
* **变量命名**：输入中文或描述，自动生成符合规范的变量或类名。
* **单词查询**：详细的单词解释及用法示例。
* **灵活配置**：支持自定义 LLM API 地址和密钥。
* **支持管道**: 支持 stdin 作为输入流，可以将其他命令输出结果直接传入进行翻译。

## 使用方法

### 设置模型提供商

首次使用前，请设置你的模型提供商 API 地址和 Key：

```bash
trans --auth "https://api.openai.com/v1,your_api_key"

```

### 常用命令

* **默认将文本翻译成英文**:
```bash
trans "你好世界"

```

* **翻译为中文**:
```bash
trans -c "Hello world"

```


* **生成变量名**:
```bash
trans -n "获取用户信息"

```


* **查单词**:
```bash
trans -w "meticulous"
trans -w "编程"

```


* **切换模型**:
```bash
trans -m "gpt-4o"

```



## 命令行选项

```text
-a, --auth strings    设置模型提供商: <base_url>,<api_key>
-c, --chinese         翻译成中文
-h, --help            查看帮助信息
-m, --model string    切换模型 (使用 'list' 查看可用模型)
-n, --name            生成变量/类名
-p, --prompt          原始提示词模式
-w, --word            单词查询模式

```


