open Unix

type item_kind = File | Dir | Phone | Error | MacBinHex | DOSBin | UnixUUENC | IdxSearch | Telnet | Bin | RServer | TN3270 | GIF | IMG | Info
type gopher_line = {kind: item_kind; content: string; selector: string; domain: string; port: int}

let kind_of_str k =
  match k with
  | "0" -> File
  | "1" -> Dir
  | "2" -> Phone
  | "3" -> Error
  | "4" -> MacBinHex
  | "5" -> DOSBin
  | "6" -> UnixUUENC
  | "7" -> IdxSearch
  | "8" -> Telnet
  | "9" -> Bin
  | "+" -> RServer
  | "T" -> TN3270
  | "g" -> GIF
  | "I" -> IMG
  | "i" -> Info
  | _ -> File

let init_socket addr port =
  let inet_addr = (gethostbyname addr).h_addr_list.(0) in
  let sockaddr = ADDR_INET (inet_addr, port) in
  let sock = socket PF_INET SOCK_STREAM 0 in
  connect sock sockaddr;

  (* file descriptor -> channel *)
  let outchan = out_channel_of_descr sock in
  let inchan = in_channel_of_descr sock in
  (inchan, outchan)

let read_page inchan =
  let lines = ref [] in
  try
    while true do
      let line = input_line inchan in
      lines := line :: !lines
    done;
    lines
  with End_of_file ->
    lines

let process_line line =
  let components = String.split_on_char '\t' line in
  if List.length components < 4 then {kind = (kind_of_str (List.nth components 0)); content = ""; selector = ""; domain = ""; port = 0}
  else {kind = (kind_of_str (List.nth components 0)); content = (List.nth components 0); selector = "test"; domain = "test"; port = 70};;

let () =
  let ic, oc = init_socket "gopher.floodgap.com" 70 in
  output_char oc '\n';
  flush oc;
  let page = read_page ic in
  List.iter ( fun l -> print_endline (process_line l).content) !page
;;