function! codesearch#test() abort
  new
  let n = winnr()
  execute n . 'wincmd w'

  setlocal buftype=acwrite
  f codesearch
  autocmd BufWriteCmd <buffer> call s:search()
endfunction

function! s:search()
  setlocal nomodified
  let body = join(getline(1, "$"), "\n")
  echo body
endfunction
