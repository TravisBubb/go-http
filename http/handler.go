package http

type Handler func(*Context)

type HandlerPipeline []Handler
