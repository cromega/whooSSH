window.onload = function () {
    var term = new Terminal();
    term.open(document.getElementById("#terminal"));
    term.writeln("whooSSH!")

    term.on("key", function (key, ev) {
        var printable = (
            !ev.altKey && !ev.altGraphKey && !ev.ctrlKey && !ev.metaKey
        );

        if (ev.keyCode == 13) {
            term.writeln("");
        } else if (ev.keyCode == 8) {
            if (term.x > 2) {
                term.write('\b \b');
            }
        } else if (printable) {
            term.write(key);
        }
    });

    var s = new WebSocket("ws://localhost:8080/whooSSH");
    s.onopen = function (event) {
        console.log("yay");
        term.attach(s);
    }

    s.onmessage = function (event) {
        console.log(event)
        term.writeln(event.data);
    }

    s.onclose = function (event) {
        term.detach(socket);
    }
}