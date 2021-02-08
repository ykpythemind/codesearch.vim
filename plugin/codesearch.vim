if exists('g:loaded_codesearch')
  finish
endif
let g:loaded_codesearch = 1

command! Hoge call codesearch#test()
