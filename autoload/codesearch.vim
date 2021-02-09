function! codesearch#open() abort
  new
  let n = winnr()
  execute n . 'wincmd w'
  set winfixheight
  resize 10

  setlocal buftype=acwrite
  setfiletype codesearch

  " 既にバッファが存在していたらそれを削除する
  let bufn = get(s:, 'codesearchbufn', -1)
  if bufn > 0
    exe 'bd!' . bufn
  endif

  let s:codesearchbufn = bufnr('')

  let pos = getpos(".")
  call append(1, "▿ includes")
  call append(2, "")
  call append(3, "▿ excludes")

  f codesearch
  autocmd BufWriteCmd <buffer> call s:search()
  autocmd FileType codesearch nnoremap <buffer> q <C-w>c
  $
  startinsert
  call setpos('.', pos)
endfunction

function! s:search()
  setlocal nomodified

  let buffer_contents = getbufline(bufnr(''), 1, '$')
  let tmpfile = tempname()

  call writefile(buffer_contents, tmpfile)

  let options = {
    \ 'stdoutbuf': '',
    \ 'num_matches': 0,
    \ }

  " let editor = has('nvim') ? s:nvim : s:vim

  call setqflist([])

  if exists('s:id')
    silent! call jobstop(s:id)
  endif

  let s:has_err = 0

  let cmd = ['codesearch-vim', tmpfile, '--cwd', getcwd()]

  let s:id = jobstart(cmd, extend(options, {
        \ 'on_stdout': function('s:on_stdout_nvim'),
        \ 'on_stderr': function('s:on_stderr_nvim'),
        \ 'stdout_buffered': 1,
        \ 'stderr_buffered': 1,
        \ 'on_exit': function('s:on_exit'),
        \ }))
endfunction

function! s:echo_err(msg) abort
  echohl ErrorMsg | echomsg a:msg | echohl None
endfunction

func! s:on_stdout_nvim(job_id, data, event) dict abort
  if !exists('s:id')
    return
  endif

  let lcandidates = []

  " https://github.com/mhinz/vim-grepper/blob/e9004ce564891412cfe433cfbb97295cccd06b39/plugin/grepper.vim#L145
  if len(a:data) > 1 || empty(a:data[-1])
    " Second-last item is the last complete line in a:data.
    let acc_line = self.stdoutbuf . a:data[0]
    let lcandidates = (empty(acc_line) ? [] : [acc_line]) + a:data[1:-2]
    let self.stdoutbuf = ''
  endif
  " Last item in a:data is an incomplete line (or empty), append to buffer
  let self.stdoutbuf .= a:data[-1]

  caddexpr lcandidates
  let self.num_matches += len(lcandidates)
endf

function! s:on_exit(job, data, event) dict abort
  unlet! s:id

  if s:has_err
    " fixme 見つからなかったときにここで弾かれる
    return
  endif

  cwindow
  if self.num_matches > 0
    let msg = printf("matched: %d", self.num_matches)
  else
    let msg = "no matches found"
  endif
  echo msg
endfunction

function! s:on_stderr_nvim(job, msg, event) dict abort
  if len(a:msg) == 1
    return
  endif

  let s:has_err = 1
  call s:echo_err(a:msg)
endfunction
