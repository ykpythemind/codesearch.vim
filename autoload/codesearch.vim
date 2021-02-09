function! codesearch#open() abort
  new
  let n = winnr()
  execute n . 'wincmd w'

  setlocal buftype=acwrite

  " 既にバッファが存在していたら上書きする
  let bufn = get(s:, 'codesearchbufn', -1)
  if bufn > 0
    exe 'bd!' . bufn
  endif

  let s:codesearchbufn = bufnr('')

  f codesearch
  " set noconfirm
  autocmd BufWriteCmd <buffer> call s:search()
  $
  startinsert
endfunction

function! s:search()
  setlocal nomodified
  let body = join(getline(1, "$"), "\n")
  echo body
endfunction
