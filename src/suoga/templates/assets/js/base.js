function make_clouds() {
    var t = [{
        src: "/assets/images/cloud-0.png",
        className: "cloud back",
        style: {
            width: "254px",
            height: "159px",
            left: .4150208376753395,
            top: .8482790827243787,
            transform: "scaleX(-1)"
        }
    }, {
        src: "/assets/images/cloud-1.png",
        className: "cloud front",
        style: {
            width: "231px",
            height: "117px",
            left: .08151647217334701,
            top: .46384778619445943,
            transform: ""
        }
    }, {
        src: "/assets/images/cloud-2.png",
        className: "cloud front",
        style: {
            width: "66px",
            height: "37px",
            left: .748033557779848,
            top: .22765147586875178,
            transform: ""
        }
    }, {
        src: "/assets/images/cloud-0.png",
        className: "cloud front",
        style: {
            width: "114px",
            height: "71px",
            left: .9580076354609097,
            top: .5181917598421091,
            transform: ""
        }
    }, {
        src: "/assets/images/cloud-0.png",
        className: "cloud back",
        style: {
            width: "96px",
            height: "60px",
            left: .526598813402908,
            top: .828749451839631,
            transform: "scaleX(-1)"
        }
    }, {
        src: "/assets/images/cloud-0.png",
        className: "cloud front",
        style: {
            width: "72px",
            height: "45px",
            left: .43174032452284195,
            top: .03627323642266411,
            transform: "scaleX(-1)"
        }
    }, {
        src: "/assets/images/cloud-0.png",
        className: "cloud back",
        style: {
            width: "84px",
            height: "53px",
            left: .9296373513977365,
            top: .2143312531352375,
            transform: ""
        }
    }, {
        src: "/assets/images/cloud-1.png",
        className: "cloud front",
        style: {
            width: "157px",
            height: "79px",
            left: .8394192676334562,
            top: .06256812600484052,
            transform: "scaleX(-1)"
        }
    }, {
        src: "/assets/images/cloud-1.png",
        className: "cloud front",
        style: {
            width: "129px",
            height: "66px",
            left: .5289903611035771,
            top: .44941927870774023,
            transform: ""
        }
    }, {
        src: "/assets/images/cloud-2.png",
        className: "cloud front",
        style: {
            width: "191px",
            height: "107px",
            left: .5054580108916613,
            top: .21665631039514555,
            transform: "scaleX(-1)"
        }
    }, {
        src: "/assets/images/cloud-1.png",
        className: "cloud front",
        style: {
            width: "257px",
            height: "130px",
            left: .711964549651326,
            top: .9866528842991085,
            transform: "scaleX(-1)"
        }
    }, {
        src: "/assets/images/cloud-0.png",
        className: "cloud front",
        style: {
            width: "160px",
            height: "100px",
            left: .8804341424789892,
            top: .9525512115988461,
            transform: "scaleX(-1)"
        }
    }, {
        src: "/assets/images/cloud-0.png",
        className: "cloud front",
        style: {
            width: "189px",
            height: "118px",
            left: .11523417887783305,
            top: .21620306890331475,
            transform: ""
        }
    }, {
        src: "/assets/images/cloud-2.png",
        className: "cloud back",
        style: {
            width: "167px",
            height: "93px",
            left: .5745663098156899,
            top: .40474044003106946,
            transform: "scaleX(-1)"
        }
    }, {
        src: "/assets/images/cloud-2.png",
        className: "cloud back",
        style: {
            width: "211px",
            height: "118px",
            left: .640291368531211,
            top: .854708255363859,
            transform: "scaleX(-1)"
        }
    }, {
        src: "/assets/images/cloud-2.png",
        className: "cloud back",
        style: {
            width: "228px",
            height: "128px",
            left: .9868028690238078,
            top: .3390108865793462,
            transform: "scaleX(-1)"
        }
    }];
    p = document.getElementsByTagName("header")[0];
    e = document.getElementById("clouds");
    s = Math.min(p.offsetHeight / 750 * t.length, t.length);
    console.log(s, "dd")
    for (var l = 0; l < s; l++) {
        console.log("swewe")
        var c = t[l],
            a = document.createElement("img");
        console.log(a, "s")
        Object.assign(a, c),
            Object.assign(a.style, c.style);
        var r = p.offsetWidth;
        a.style.left = r * c.style.left + "px";
        var o = p.offsetHeight;
        a.style.top = o * c.style.top + "px",
            e.appendChild(a)
        console.log(a, "s23")
    }
    console.log("122")
    var n = 0,
        i = 0,
        d = 14,
        m = 5,
        f = e.children;

    function h(t, e, s, l) {
        var c = (e - s) / 1e3 * l;
        if (c < 1)
            return !1;
        for (var a = 0; a < f.length; a++) {
            var r = f[a];
            if (r.className.includes(t)) {
                var o = parseInt(r.style.left) - c,
                    n = parseInt(r.style.width);
                if (o < -n) {
                    var i = p.offsetWidth;
                    o = i - -(o + n) % i
                }
                r.style.left = o + "px"
            }
        }
        return !0
    }
    window.requestAnimationFrame(function t(e) {
        e - Math.min(i, n) < 50 || (h("front", e, n, d) && (n = e),
                h("back", e, i, m) && (i = e)),
            window.requestAnimationFrame(t)
    })
}