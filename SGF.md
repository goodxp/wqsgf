# SGF File Format FF[4] Notes

This doc is basically my notes taken during the read on SGF standard, 
from : https://www.red-bean.com/sgf/index.html
It is incomplete but handy for coding and quick reference. 

* A node is the smallest unit visible to the user. 
For example: if you step forward and see a new move on the board and 
a comment in the comment window plus some markup on the board then all 
this information is represented by different properties which are parts 
of the same node.
* The first branch (variation A or 1) is the main branch.
* The last node of the game should contain the last move on the board. 

## EBNF
* Collection = GameTree { GameTree }
* GameTree   = "(" Sequence { GameTree } ")"
* Sequence   = Node { Node }
* Node       = ";" { Property }  //Only one of each property is allowed per node
* Property   = PropIdent PropValue { PropValue }
* PropIdent  = UcLetter { UcLetter }
* PropValue  = "[" CValueType "]"
* CValueType = (ValueType | Compose)
* ValueType  = (None | Number | Real | Double | Color | SimpleText |
        Text | Point  | Move | Stone)

* 'list of':    PropValue { PropValue }
* 'elist of':   ((PropValue { PropValue }) | None)
            In other words elist is list or "[]".

## Propety Types
* root: global game 'attributes'. may only appear in root nodes.
* game-info: info about the game played. 
    usually stored in root nodes.
    may be only one game-info node on any path within a game tree.
* setup: current position.
* move: move made.
* no type: appear anywhere.

* inherit: affect the node containing it and subsequent child nodes.
    until a new setting or cleared. E.g. variation tree: VM...VM[]

* private: user-defined. escape "]" => "\]"

## Property Value Types
* UcLetter   = "A".."Z"
* Digit      = "0".."9"
* None       = ""

* Number     = [("+"|"-")] Digit { Digit }
* Real       = Number ["." Digit { Digit }]

* Double     = ("1" | "2") 
    //for annotation properties (normal | emphasized)
* Color      = ("B" | "W")

* SimpleText = { any character (handling see below) }
    //no newline! ...as by Text
* Text       = { any character (handling see below) } 
    // White spaces other than linebreaks are converted to space (e.g. no tab, vertical tab, ..).
    // soft linebreaks preceded by a "\" (converted to "")
    // Escaping: "\" is the escape character for "]", "\" and ":".

* Point      = game-specific
* Move       = game-specific
* Stone      = game-specific
// in Go(GM[1]), Point=Stone. e.g. "aa". 19x19=a..s, 52x52=a..zA..Z
// [tt] = pass
// "list of point" or "elist of point" may be compressed. use the upper left and lower right corner of the rectangle:
// (point | composition of point ":" point)
// e.g.	(;GM[1]SZ[9]FF[4]
            AB[ac:ic]AW[ae:ie]
            ;DD[aa:bi][ca:ce]
            VW[aa:bi][ca:ee][ge:ie][ga:ia][gc:ic][gb][ib])

* Compose    = ValueType ":" ValueType

## Properties (67+5=72)
* Move(4): 			    B, KO, MN, W
* Setup(4): 			AB, AE, AW, PL
* Node Annotation(8): 	C, DM, GB, GW, HO, N, UC, V
* Move Annotation(4): 	BM, DO, IT, TE
* Markup(9): 			AR, CR, DD, LB, LN, MA, SL, SQ, TR
* Root(6): 			    AP, CA, FF, GM, ST, SZ
* Game Info(21): 		AN, BR, BT, CP, DT, EV, GN, GC, ON, OT, PB, 
                        PC, PW, RE, RO, RU, SO, TM, US, WR, WT
* Timing(4): 			BL, OB, OW, WL
* Miscellaneous(3): 	FG, PM, VW
* Go(GM[1])(4):		    HA, KM, TB, TW

* Non-Standard(5):		BC, WC, LT, LC, GK

## Prop examples:
### Move:
* B[ab] W[cd]
* KO //invalid move in node, skip rule check
* MN[26]B[ab] //move number

### Setup:
* AB[aa][bb] AW[cc] //Add black/white stones to the board
* AE[aa][bb] //Clear stones on the board
* PL[B] //whose turn it is to play

### Node Anno:
* C[this is a comment]
* DM[1] DM[2] // the position is even. Not with UC, GB or GW
* GB[1] GW[1] // good for black/white
* HO[1] //hotspot, interesting
* N[funny node] //provides a name for the node
* UC[1] //position unclear
* V[26.5] V[-26.5] //value for the node. In Go, estimated score, positive=good for black

### Move Anno:
These are exclusive, NOT used within the same node. 
Shall be with B[] W[].
* BM[1] //bad move
* DO //doubtful
* IT //interesting
* TE[1] //tesuji (good move)

### Markup:
* AR[fg] AR[aa:cc][aa:dd] // arrow(s)
* CR[aa] CR[aa][bb] // circle(s)
* DD[aa][bb] DD[] // Dim (grey out), undim everything
* VM[aa:this is good][bb:that is bad] // text shown around the point
* LN[aa][cc] // line(s)
* MA[aa][cc] // 'x'
* SL[aa][cc] // selected. inverts color?
* SQ[aa][cc] // square
* TR[aa][cc] // triangle

### Root:
* AP[CGoban:1.6.2], Unix
    [Hibiscus:2.1], Windows 95
    [IGS:5.0], Windows 95
    [Many Faces of Go:10.0], Windows 95
    [MGT:?], DOS/Unix
    [NNGS:?], Unix
    [Primiview:3.0], Amiga OS3.0
    [SGB:?], Macintosh
    [SmartGo:1.0], Windows
    //name and version number of the application to create this gametree
* CA[ISO-8859-1] CA[UTF-8] //charset
* FF[1] .. FF[4] //SGF file format. 
* GM[1] .. GM[20] //game type. Go=1(default). Othello = 2, chess = 3,
    Gomoku+Renju = 4, Nine Men's Morris = 5, Backgammon = 6,
    Chinese chess = 7, Shogi = 8, Lines of Action = 9,
    Ataxx = 10, Hex = 11, Jungle = 12, Neutron = 13,
    Philosopher's Football = 14, Quadrature = 15, Trax = 16,
    Tantrix = 17, Amazons = 18, Octi = 19, Gess = 20.
* ST[0] .. ST[3] //how variations should be shown
    //show variations of successor node (children) (value: 0)
    //show variations of current node   (siblings) (value: 1)
    //no (auto-) board markup (value: 2)
* SZ[19] SZ[12:8] //board size, 19x19, 12x8

### Game info:
* AN[Lee Wong] //name of who made the game annotations
* BR[9d] BR[1k*] BR[?] //rank of the black player
* BT[Milk Cow] // name of the black team
* CP[copyright 2019 goodxp for annotations] //copyright
* DT[2001-12-06,07,08] // when the game was played.
* EV[Milk Cow World Tournament] // event name
* GN[All Stones Killed] // game name
* GC[only milk cows allowed...] // extra game info
* ON[Chinese fuseki] //opening played info
* OT[25 moves / 10 min] //method used for overtime
* PB[Shabi Xi] PW[Erhuo Mao] //black player, white player
* PC[Restroom] //where the game was played
* RE[0] RE[B+10.5] RE[W+R] RE[B+T] RE[W+F] RE[] RE[?]
    // result of game. 
    // draw 0, black win 10.5, White resigned, black OT, forfeit, none, unknown
* RO[World Cup (final)] //round-number
* RU[AGA] RU[Chinese] RU[Japanese] //game rules
* SO[CNN] //name of the source
* TM[120] //time limits of game (seconds)
* US[Trump Bush] //name of user who entered the game
* WR[9d] //white player rank
* WT[No Milk Cow] //name of white team

### Timing:
* BL[10] //time(seconds) left for black after the move made
* OB[9] //number of black moves left to play in this byo-yomi period
* OW[8] //number of white moves left to play in this byo-yomi period
* WL[30] //time(seconds) left for white after the move made

### Misc:
These are basically for printing...
* FG FG[1:xxx][2:xx] // figure start: divide a game into different figures
* PM[0] // for printing, 0=no move numbers...
* VM VM[aa:cc] // View only part of the board

### Go(GM[1]):
* TB[aa][bb] //black territory or area
* TW[aa][bb] //white territory or area
* KM[7.5] //komi
* HA[2] //handicap

### Non-Standard:
* BC[kr] WC[cn] //country of black and white players
* LT[]LC[5]GK[1] // unknown

## SGF Examples

### A short SGF
    (;FF[4]C[root](;C[a];C[b](;C[c])(;C[d];C[e]))(;C[f](;C[g];C[h];C[i])(;C[j])))

To simplify the above model:

    (root(ab(c)(de))(f(ghi)(j)))

The corresponding tree structure: 
    // horizontal line (-) represents a variation gametree.

		(root
		 |
		(a ------- (f
		 |          |
		 b         (g - (j)))
		 |          |
		(c) - (d    h
			   |    |
			   e))  i)

### A real world SGF file content
From http://gokifu.com/index.php

    (;FF[4]GM[1]SZ[19]CA[UTF-8]SO[gokifu.com]BC[kr]WC[cn]EV[]PB[朴廷桓]BR[九段]PW[柯洁]WR[九段]KM[7.5]DT[2019-10-13]RE[白胜黑]TM[120]LT[]LC[5]GK[1];B[pd];W[dd];B[pp];W[dp];B[cq];W[dq];B[cp];W[co];B[bo];W[bn];B[cn];W[do];B[bm];W[bp];B[an];W[bq];B[fc];W[cf];B[qn];W[nc];B[qf];W[ec];B[oc];W[nd];B[eb];W[db];B[ed];W[dc];B[jc];W[ng];B[ph];W[or];B[pr];W[pq];B[oq];W[qq];B[qr];W[op];B[nq];W[qp];B[po];W[nr];B[rr];W[mq];B[np];W[lp];B[no];W[ip];B[ni];W[kf];B[lc];W[kd];B[ld];W[jd];B[kc];W[id];B[lf];W[lg];B[ke];W[jf];B[gb];W[mf];B[le];W[mj];B[of];W[og];B[pg];W[oe];B[pf];W[me];B[mh];W[mg];B[nj];W[ml];B[jq];W[jp];B[kq];W[kp];B[mk];W[lk];B[nk];W[pe];B[qd];W[on];B[ll];W[oo];B[nn];W[om];B[qo];W[nm];B[mn];W[lm];B[kl];W[mm];B[ln];W[km];B[kk];W[ki];B[fr];W[er];B[lr];W[mr];B[iq];W[gr];B[kh];W[gq];B[jh];W[kg];B[if];W[ig];B[je];W[lh];B[li];W[jj];B[ii];W[jg];B[hf];W[jl];B[kj];W[ji];B[jk];W[ih];B[jm];W[kn];B[pl];W[gd];B[fd];W[ge];B[ie];W[gf];B[hg];W[hh];B[gg];W[ef];B[cj];W[ea];B[gc];W[il];B[ik];W[hk];B[hl];W[im];B[ij];W[hj];B[hm];W[in];B[hi];W[gi];B[lj];W[gh];B[gj];W[fg];B[gk];W[hd];B[mb];W[nb];B[fb];W[bi];B[ci];W[bh];B[de];W[ce];B[dg];W[dh];B[ch];W[cg];B[eh];W[fh];B[fe];W[ff];B[df];W[ee];B[di];W[ka];B[cd];W[bd];B[na];W[oa];B[ma];W[mc];B[lb];W[ib];B[ic];W[hc];B[hb];W[ob];B[jb];W[qe];B[re];W[fa];B[ia];W[ga])
