# Gher
Short for `Generic Http Encoding Responder`, gher is a simple middleware that uses generics to
automatically encode and respond to http requests. It is useful when you don't care to handle the
underlying http responses and just want to write functions that take in some generic type and
return some generic type or an error.