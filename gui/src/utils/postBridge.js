const postBridge = data => fetch('/v2/bridge_types', {
  method: 'POST',
  body: data,
  headers: {
    Authorization: 'Basic Y2hhaW5saW5rOnR3b2NoYWlucw==',
    'Content-Type': 'application/json'
  }
})

export default postBridge
