TODO:
- make initial ask
- make request
- decode response
- if response contains toolCalls OR the body is valid JSON OR is JSON wrapped in xml tags of <tool_calls>
  then do a toolCall loop
    - make tool calls UNTIL
      no toolCalls AND the body is not valid JSON OR is JSON wrapped in xml tags of <tool_calls>