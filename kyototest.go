/*
Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"zylo/reiwa"
)

//go:embed kyototest.dat
var cityMultiList string

var mul2map map[string]int

var mul2sum int

func init() {
	reiwa.CityMultiList = cityMultiList
	reiwa.OnAssignEvent = onAssignEvent
	reiwa.OnInsertEvent = onInsertEvent
	reiwa.OnDeleteEvent = onDeleteEvent
	reiwa.OnAcceptEvent = onAcceptEvent
	reiwa.OnPointsEvent = onPointsEvent
	reiwa.AllowBandRange(reiwa.K1900, reiwa.M5600)
	reiwa.AllowModeRange(reiwa.CW, reiwa.AM)
	reiwa.AllowRcvd(`^([A-Z]{1}\d{2}|[A-Z]{2})(\d{3}|[A-Z]{2})$`)
}

func onAssignEvent(contest, configs string) {
	mul2map = make(map[string]int)
	mul2sum = 0
}

func onInsertEvent(qso *reiwa.QSO) {
        if len(qso.GetMul2()) >= 3 {
	   _, ok := mul2map[qso.GetMul2()]
	   if ok {
	      mul2map[qso.GetMul2()] += 1
	   } else {
	     mul2sum += 1
	     mul2map[qso.GetMul2()] = 1
	     qso.SetNote("new multi " + qso.GetMul2())
	   }
	}
}

func onDeleteEvent(qso *reiwa.QSO) {
        if len(qso.GetMul2()) >= 3 {
	   mul2map[qso.GetMul2()] -= 1
	   if mul2map[qso.GetMul2()] == 0 {
		mul2sum -= 1
		delete(mul2map, qso.GetMul2())
	   }
	}
}

func isInPref(mul1 string) bool {
	return len(mul1) >= 3
}

func isIn_to_In(rcvd_mul1, sent_mul1 string) bool {
	return isInPref(rcvd_mul1) == true && isInPref(sent_mul1) == true
}

func isOut_to_Out(rcvd_mul1, sent_mul1 string) bool {
	return isInPref(rcvd_mul1) == false && isInPref(sent_mul1) == false
}

func score(rcvd_mul1, sent_mul1 string) byte {
	switch {
	case isIn_to_In(rcvd_mul1, sent_mul1):
		return 2
	case isOut_to_Out(rcvd_mul1, sent_mul1):
		return 0
	default:
		return 1
	}
}

func onAcceptEvent(qso *reiwa.QSO) {
	rcvd := qso.GetRcvdGroups()
	sent := qso.GetSentGroups()
	rcvd_mul1, rcvd_mul2 := rcvd[1], rcvd[2]
	sent_mul1 := sent[1]
	qso.Score = score(rcvd_mul1, sent_mul1)
	qso.SetMul1(rcvd_mul1)
	qso.SetMul1(rcvd_mul2)
}

func onPointsEvent(score, mults int) int {
	return score * (mults + mul2sum)
}
