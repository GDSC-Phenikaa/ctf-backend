package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/GDSC-Phenikaa/ctf-backend/db"
	"github.com/GDSC-Phenikaa/ctf-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type SeedQuestion struct {
	Content       string
	Type          string
	Options       string
	CorrectAnswer string
	Points        int
}

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	database, err := db.Connect()
	if err != nil {
		panic(err)
	}

	fmt.Println("Deploying 100 ACTUAL cybersecurity questions...")

	advancedModule := models.Module{Title: "The Master Challenge Server", Description: "100 varied realistic cyber security questions.", Order: 10}
	if err := database.Where("title = ?", advancedModule.Title).First(&advancedModule).Error; err != nil {
		database.Create(&advancedModule)
	}

	lesson := models.Lesson{ModuleID: advancedModule.ID, Title: "The Ultimate Knowledge Base", Content: "This lesson contains 100 real-world cyber questions.", Order: 1}
	if err := database.Where("title = ?", lesson.Title).First(&lesson).Error; err != nil {
		database.Create(&lesson)
	}

	realQuestions := []SeedQuestion{
		// Networking
		{"What port does standard unencrypted FTP operate on?", "mcq", `["21", "22", "23", "25"]`, "21", 50},
		{"Which protocol uses port 22 by default?", "mcq", `["Telnet", "SSH", "HTTP", "SMB"]`, "SSH", 50},
		{"What is port 3389 commonly used for?", "mcq", `["DNS", "RDP", "LDAP", "VNC"]`, "RDP", 50},
		{"Which OSI layer determines the routing of packets?", "mcq", `["Layer 2", "Layer 3", "Layer 4", "Layer 7"]`, "Layer 3", 50},
		{"What is the default port for DNS?", "text", `[]`, "53", 50},
		{"What is the default port for SMB?", "mcq", `["139", "445", "both", "neither"]`, "both", 75},
		{"Which utility is a comprehensive network and port scanner?", "mcq", `["Nmap", "Wireshark", "Burp Suite", "Hashcat"]`, "Nmap", 50},
		{"Which protocol ensures secure delivery over HTTP?", "mcq", `["UDP", "TLS", "IPsec", "SNMP"]`, "TLS", 50},
		{"In a TCP Handshake, what packet follows SYN?", "mcq", `["ACK", "FIN", "SYN-ACK", "RST"]`, "SYN-ACK", 50},
		{"What does a subnet mask of 255.255.255.0 correspond to in CIDR notation?", "text", `[]`, "/24", 75},
		{"Which port is commonly used for LDAP over SSL (LDAPS)?", "text", `[]`, "636", 100},
		{"What does ARP stand for?", "mcq", `["Address Resolution Protocol", "Active Record Protocol", "Advanced Routing Protocol", "Asymmetric RSA Protocol"]`, "Address Resolution Protocol", 50},
		{"Which type of scan does nmap execute with the '-sS' flag?", "mcq", `["Connect scan", "SYN stealth scan", "UDP scan", "NULL scan"]`, "SYN stealth scan", 100},
		{"What is a MAC address?", "mcq", `["A 32-bit logical address", "A 48-bit physical address", "A 128-bit physical address", "A cryptographic hash"]`, "A 48-bit physical address", 50},
		{"What Protocol Data Unit (PDU) is used at the Transport layer?", "mcq", `["Frame", "Packet", "Segment", "Data"]`, "Segment", 50},
		
		// Cryptography
		{"Which cryptographic algorithm is asymmetric?", "mcq", `["AES", "DES", "RSA", "Blowfish"]`, "RSA", 75},
		{"What is the maximum key size of AES?", "text", `[]`, "256", 50},
		{"Which hash algorithm produces a 128-bit digest and is severely broken?", "mcq", `["SHA-1", "MD5", "SHA-256", "Bcrypt"]`, "MD5", 50},
		{"What attack takes advantage of two different inputs producing the same hash value?", "mcq", `["Birthday Attack", "Rainbow Table Attack", "Brute Force Attack", "Man-in-the-Middle Attack"]`, "Birthday Attack", 100},
		{"In PKI, what is CSR an acronym for?", "mcq", `["Certificate Signing Request", "Common Security Rule", "Certificate Security Response", "Common Signature Request"]`, "Certificate Signing Request", 75},
		{"What is the mathematical concept underlying the security of RSA?", "mcq", `["Elliptic Curves", "Discrete Logarithms", "Prime Number Factorization", "Symmetric Substitution"]`, "Prime Number Factorization", 100},
		{"What length is a standard SHA-256 hash in hex characters?", "text", `[]`, "64", 75},
		{"What mode of operation for block ciphers uses a nonce combined with a counter to encrypt data?", "mcq", `["ECB", "CBC", "CTR", "GCM"]`, "CTR", 100},
		{"Why is ECB considered an insecure block cipher mode?", "mcq", `["It is too slow", "It requires too much padding", "Identical plaintext blocks produce identical ciphertext blocks", "It cannot be implemented in hardware"]`, "Identical plaintext blocks produce identical ciphertext blocks", 50},
		{"What is the purpose of a 'salt' in password hashing?", "mcq", `["To speed up verification", "To prevent rainbow table attacks", "To enable symmetric decryption of passwords", "To bypass collision resistance"]`, "To prevent rainbow table attacks", 50},
		{"Diffie-Hellman securely accomplishes what cryptographic objective?", "mcq", `["Hashing", "Digital Signatures", "Key Exchange over an insecure channel", "Symmetric block encryption"]`, "Key Exchange over an insecure channel", 100},
		{"What relies on the difficulty of finding the discrete logarithm of a random elliptic curve element?", "mcq", `["RSA", "ECC", "AES", "DES"]`, "ECC", 150},
		{"Which encoding mechanism translates binary into a string composed of A-Z, a-z, 0-9, +, and /?", "mcq", `["Hexadecimal", "Base32", "Base64", "URL Encoding"]`, "Base64", 50},
		{"What are the last two characters often used for padding in Base64 strings?", "text", `[]`, "==", 50},
		{"Vigenère cipher is an example of what type of cipher?", "mcq", `["Polyalphabetic substitution", "Monoalphabetic substitution", "Transposition", "Stream cipher"]`, "Polyalphabetic substitution", 100},
		
		// Web Exploitation (OWASP)
		{"What is SQL Injection (SQLi)?", "mcq", `["Modifying an HTML structure", "Exploiting bad database queries", "Cross-site request passing", "A buffer overflow"]`, "Exploiting bad database queries", 50},
		{"What does a classic authentication bypass SQL injection look like?", "text", `[]`, "' OR '1'='1", 100},
		{"Which type of XSS occurs when a payload is permanently saved on the target application?", "mcq", `["Reflected XSS", "Stored XSS", "DOM XSS", "Blind XSS"]`, "Stored XSS", 50},
		{"Which HTTP header helps mitigate Cross-Site Scripting (XSS) by defining approved sources of content?", "mcq", `["Content-Security-Policy (CSP)", "X-Frame-Options", "Strict-Transport-Security", "Access-Control-Allow-Origin"]`, "Content-Security-Policy (CSP)", 100},
		{"What attack forces a logged-in victim's browser to execute an unwanted action on a web application?", "mcq", `["XSS", "SQLi", "CSRF", "LFI"]`, "CSRF", 75},
		{"What is the most common defense against CSRF attacks?", "mcq", `["Anti-CSRF Tokens", "Input Validation", "HTTPS", "WAF"]`, "Anti-CSRF Tokens", 75},
		{"What vulnerability allows an attacker to include files on a server through the web browser?", "mcq", `["SSRF", "LFI", "RCE", "XXE"]`, "LFI", 50},
		{"What payload is commonly used to access the passwd file via LFI on Linux?", "text", `[]`, "../../../../../etc/passwd", 100},
		{"If an application processes user input by fetching a URL internally, what vulnerability might exist?", "mcq", `["SSRF", "LFI", "SQLi", "Open Redirect"]`, "SSRF", 100},
		{"What XML entity vulnerability can result in reading arbitrary files or SSRF?", "mcq", `["XPath Injection", "XXE", "XML Bomb", "XSLT Injection"]`, "XXE", 100},
		{"What does the acronym CORS stand for?", "mcq", `["Cross-Origin Resource Sharing", "Cross-Origin Remote Scripting", "Central Origin Redirect System", "Cross-Object Reference Source"]`, "Cross-Origin Resource Sharing", 50},
		{"Which HTTP method is specifically designed to perform preflight checks in CORS?", "mcq", `["PUT", "TRACE", "OPTIONS", "HEAD"]`, "OPTIONS", 100},
		{"What is IDOR?", "mcq", `["Insecure Direct Object Reference", "Internal Database Operational Request", "Inline Direct Object Response", "Incomplete Domain Object Routing"]`, "Insecure Direct Object Reference", 75},
		{"What is the root cause of a directory traversal attack?", "mcq", `["Weak passwords", "Missing input validation of filepath operations", "Lack of HTTPS", "Insecure hashing"]`, "Missing input validation of filepath operations", 50},
		{"A JSON Web Token (JWT) typically consists of how many parts separated by dots?", "text", `[]`, "3", 50},
		{"Which part of a JWT contains the cryptographic mechanism asserting its authenticity?", "mcq", `["Header", "Payload", "Signature", "Footer"]`, "Signature", 50},
		{"What happens when a JWT header specifies 'alg':'none' and the backend blindly trusts it?", "mcq", `["The token becomes encrypted", "Authentication bypass occurs as signature validation is skipped", "The token is invalidated immediately", "The payload expands infinitely"]`, "Authentication bypass occurs as signature validation is skipped", 150},
		{"Which tag allows reading cookies in an XSS attack?", "text", `[]`, "document.cookie", 100},
		{"Which flag ensures that a cookie cannot be accessed via JavaScript?", "mcq", `["Secure", "HttpOnly", "SameSite", "Path"]`, "HttpOnly", 75},
		{"What proxy tool is widely considered the industry standard for web vulnerability hunting?", "mcq", `["Wireshark", "Burp Suite", "Metasploit", "Nmap"]`, "Burp Suite", 50},

		// Reverse Engineering / Binary Exploitation (Pwn)
		{"In a buffer overflow, what crucial register or value is overwritten to hijack execution flow?", "mcq", `["EAX / RAX", "Instruction Pointer (EIP / RIP) / Return Address", "EBX / RBX", "Stack Base Pointer (EBP)"]`, "Instruction Pointer (EIP / RIP) / Return Address", 100},
		{"What mitigates buffer overflows by randomizing the memory locations of key data areas?", "mcq", `["NX / DEP", "Stack Canaries", "ASLR", "RELRO"]`, "ASLR", 100},
		{"What does NX (No-eXecute) / DEP do?", "mcq", `["Randomizes memory addresses", "Encrypts the stack", "Prevents code execution from the stack/heap", "Detects buffer overflows before they happen"]`, "Prevents code execution from the stack/heap", 100},
		{"What is a Stack Canary?", "mcq", `["A random value placed on the stack before the return pointer", "A bird used in coal mines", "A type of firewall", "A specialized debugger command"]`, "A random value placed on the stack before the return pointer", 100},
		{"Which class of vulnerability uses '%x', '%s', and '%n' maliciously?", "mcq", `["Integer Overflow", "Format String Vulnerability", "Heap Use-After-Free", "Race Condition"]`, "Format String Vulnerability", 100},
		{"What does ROP stand for in binary exploitation?", "mcq", `["Return Oriented Programming", "Remote OS Pwning", "Runtime Object Parsing", "Reverse Overhead Protocol"]`, "Return Oriented Programming", 150},
		{"In x86 assembly, what instruction is typically used to move the Stack Pointer back to the Base Pointer?", "mcq", `["RET", "LEAVE", "PUSH", "POP"]`, "LEAVE", 150},
		{"What is 'malloc' used for in C?", "mcq", `["Thread synchronization", "Dynamic memory allocation on the heap", "Opening network sockets", "Reading from standard input"]`, "Dynamic memory allocation on the heap", 50},
		{"What heap vulnerability occurs when a program frees memory but continues to use the pointer?", "mcq", `["Double Free", "Use-After-Free", "Heap Overflow", "Memory Leak"]`, "Use-After-Free", 125},
		{"What file contains the dynamic linker resolution entries on ELF binaries?", "mcq", `["GOT (Global Offset Table)", "BSS", "TEXT", "DATA"]`, "GOT (Global Offset Table)", 200},
		{"What is the first command normally run in GDB to start a program?", "text", `[]`, "run", 50},
		{"If I want to find the flag in a stripped binary, what reverse engineering framework developed by the NSA could I use?", "mcq", `["Ghidra", "IDA Pro", "Radare2", "All of the above"]`, "All of the above", 75},
		{"What is the magic number representing an ELF file header in hexadecimal?", "text", `[]`, "7f454c46", 200},
		{"What does the 'NOP Sled' do?", "mcq", `["Sliding execution safely until hitting shellcode", "Crashing the program instantly", "Accelerating processing speed", "Nullifying the return address"]`, "Sliding execution safely until hitting shellcode", 100},
		{"In 64-bit Linux (x86_64) calling convention, which register holds the first integer argument?", "mcq", `["RDI", "RSI", "RDX", "RCX"]`, "RDI", 150},

		// Forensics and OSINT
		{"What tool is commonly used to find hidden files embedded inside image files?", "mcq", `["steghide", "strings", "binwalk", "volatility"]`, "steghide", 75},
		{"What command extracts printable text from a compiled binary or data file?", "text", `[]`, "strings", 50},
		{"If you have an unknown file, what Linux command will attempt to determine its type based on magic bytes?", "text", `[]`, "file", 50},
		{"Which tool analyzes network packet captures (.pcap/pcapng files) visually?", "mcq", `["Wireshark", "Tcpdump", "Nmap", "Netcat"]`, "Wireshark", 50},
		{"What powerful Python framework extracts artifacts from RAM memory dumps?", "mcq", `["Volatility", "Autopsy", "Ghidra", "Binwalk"]`, "Volatility", 100},
		{"What protocol allows fetching domain registration details (often used in OSINT)?", "text", `[]`, "whois", 50},
		{"What search engine is designed specifically to find internet-connected devices, webcams, and open ports?", "mcq", `["Google", "Shodan", "DuckDuckGo", "Censys"]`, "Shodan", 75},
		{"What is 'Binwalk' most frequently used for?", "mcq", `["Scanning for open ports", "Extracting firmware images and file systems from binary blobs", "Cracking passwords", "Analyzing network traffic"]`, "Extracting firmware images and file systems from binary blobs", 100},
		{"If a file has the magic bytes '89 50 4E 47 0D 0A 1A 0A', what type of file is it?", "mcq", `["JPEG", "PNG", "GIF", "PDF"]`, "PNG", 100},
		{"If a file starts with '%PDF', what type of file is it?", "text", `[]`, "PDF", 50},
		{"What command can extract data from a zip file?", "text", `[]`, "unzip", 50},
		{"In Windows forensics, what central hierarchical database stores configurations and state information?", "mcq", `["Active Directory", "NTFS MFT", "The Registry", "Event Logs"]`, "The Registry", 75},
		{"Which file in Linux contains hashed passwords for all accounts?", "mcq", `["/etc/passwd", "/etc/shadow", "/etc/hash", "/etc/sudoers"]`, "/etc/shadow", 50},
		{"What is 'LSASS' responsible for in Windows?", "mcq", `["Local Security Authority Subsystem Service (handles logins & credentials)", "Linux Subsystem Architecture Service", "Low-level System Application Server", "Local Storage Access Server"]`, "Local Security Authority Subsystem Service (handles logins & credentials)", 120},
		{"What tool famously dumps plaintext credentials and hashes from LSASS memory?", "mcq", `["Metasploit", "Mimikatz", "John the Ripper", "Hydra"]`, "Mimikatz", 100},

		// Miscellaneous & Command Line
		{"What does the 'chmod' command do in Linux?", "mcq", `["Changes owner of a file", "Changes permissions of a file", "Changes the password of a user", "Creates a module"]`, "Changes permissions of a file", 50},
		{"If a file has permissions '777', what does it mean?", "mcq", `["Only root can read/write", "Owner can read, group can write, others execute", "Everyone can read, write, and execute", "It's a hidden directory"]`, "Everyone can read, write, and execute", 50},
		{"What Linux command lists all heavily detailed running processes?", "text", `[]`, "ps aux", 75},
		{"What is standard input (stdin) file descriptor number in linux?", "text", `[]`, "0", 100},
		{"What symbol is used to redirect standard output to a file, overwriting the file?", "text", `[]`, ">", 50},
		{"Which command line utility allows sending raw TCP/UDP data and is known as the 'Swiss Army Knife' of networking?", "mcq", `["curl", "wget", "netcat (nc)", "ping"]`, "netcat (nc)", 75},
		{"What payload is used to spawn an interactive bash shell in Python?", "text", `[]`, "import pty; pty.spawn('/bin/bash')", 150},
		{"What privilege escalation vulnerability exploits programs with high privileges running on behalf of an unprivileged user?", "mcq", `["Cron Job vulnerabilities", "SUID binary exploitation", "World writable files", "All of the above"]`, "All of the above", 75},
		{"What does GTFOBins provide?", "mcq", `["A list of default passwords", "A list of Linux binaries that can be exploited to bypass local security", "A compiler for Windows to Linux cross-compilation", "A database of hashes"]`, "A list of Linux binaries that can be exploited to bypass local security", 100},
		{"Which encryption standard replaced DES and remains robust today?", "text", `[]`, "AES", 50},
		{"What tool automates SQL injection attacks?", "mcq", `["Hydra", "DirBuster", "SQLMap", "Hashcat"]`, "SQLMap", 75},
		{"What command enables 'sudo' to impersonate any user, specifically checking what privileges you have?", "text", `[]`, "sudo -l", 75},
		{"A reverse shell occurs when...", "mcq", `["You connect to the victim's open listening port", "The victim's machine connects back to your listening port", "Shell access is mirrored to two attackers", "You run a command blindly via RCE"]`, "The victim's machine connects back to your listening port", 100},
		{"Which hash cracker heavily utilizes GPU acceleration?", "mcq", `["John the Ripper", "Hashcat", "Hydra", "CrackMapExec"]`, "Hashcat", 75},
		{"Which flag is widely used globally as a test or dummy domain name string?", "text", `[]`, "example.com", 50},

		// Cloud & General
		{"Which AWS service stores objects like images, files, and backups?", "mcq", `["EC2", "RDS", "S3", "Lambda"]`, "S3", 50},
		{"What is an SSRF vulnerability commonly used for when attacking an AWS EC2 instance?", "mcq", `["DDoS attacks", "Accessing the instance metadata service at 169.254.169.254", "Reading local SQL databases", "Bypassing the Web Application Firewall"]`, "Accessing the instance metadata service at 169.254.169.254", 150},
		{"In Docker, what flag runs the container with full host privileges?", "text", `[]`, "--privileged", 125},
		{"If a Docker daemon socket is exposed over TCP without TLS, what is the impact?", "mcq", `["Denial of Service", "Full Host Compromise via API endpoints", "Nothing, it requires auth by default", "Information Disclosure of image tags"]`, "Full Host Compromise via API endpoints", 150},
		{"What Kubernetes component provides cluster wide secret and state management?", "mcq", `["Kubelet", "Etcd", "Kube-proxy", "CoreDNS"]`, "Etcd", 125},
	}

	for _, rq := range realQuestions {
		q := models.Question{
			LessonID:      lesson.ID,
			Content:       rq.Content,
			Type:          rq.Type,
			Options:       rq.Options,
			CorrectAnswer: rq.CorrectAnswer,
			Points:        rq.Points,
		}

		var existing models.Question
		if err := database.Where("content = ?", q.Content).First(&existing).Error; err != nil {
			database.Create(&q)
		} else {
			q = existing
		}
	}

	// Fetch users for some random solves
	var users []models.User
	database.Where("is_admin = ?", false).Find(&users)
	if len(users) > 0 {
		rand.Seed(time.Now().UnixNano())
		// Get all questions in the advanced module
		var advQs []models.Question
		database.Where("lesson_id = ?", lesson.ID).Find(&advQs)
	
		for _, q := range advQs {
			numSolvers := rand.Intn(len(users) + 1)
			rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
			
			for j := 0; j < numSolvers; j++ {
				createLMSSolve(database, q.ID, users[j].ID, true, q.CorrectAnswer)
			}
		}
	}

	fmt.Println("Successfully deployed almost 100 absolutely real cybersecurity questions and random solves!")
}

func createLMSSolve(database *gorm.DB, questionID, userID uint, correct bool, answer string) {
	var solve models.QuestionSolve
	if err := database.Where("question_id = ? AND user_id = ?", questionID, userID).First(&solve).Error; err != nil {
		database.Create(&models.QuestionSolve{QuestionID: questionID, UserID: userID, Correct: correct, SubmittedAnswer: answer})
	}
}
