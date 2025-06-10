TODO:
- make initial ask
- make request
- decode request
- if response contains toolCalls OR the body is valid JSON
  then do a toolCall loop
    - make tool calls UNTIL
      no toolCalls AND the body is not valid JSON