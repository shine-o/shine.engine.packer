package handlers

import (
	"encoding/binary"
	"fmt"
	"github.com/go-restruct/restruct"
	"github.com/google/logger"
	"github.com/segmentio/ksuid"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ExtractCmd struct {
	Path             string
	StoreServerFiles bool
	ServerFilesPath  string
}

type GDP struct {
	Type string `struct:"[3]byte"`
	UnkData1 [5]byte
	Version string `struct:"[260]byte"`
	Num1 uint32
	NumOfFiles uint32
	UnkData3 [40]byte
	Files []File `struct-size:"NumOfFiles - 1"`
}


type File struct {
	Num0 uint64
	Name string `struct:"[264]byte"`
	Offset uint64
	Size uint64
	Num3 uint64
	UnkData4 [20]byte
	Data []byte
}

func Extract(cmd *cobra.Command, args []string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	restruct.EnableExprBeta()

	log = logger.Init("extract logger", true, true, ioutil.Discard)

	if err := cmd.Flags().Parse(args); err != nil {
		log.Fatal(err)
	}

	sourcePath, _ := cmd.Flags().GetString("source")
	destinationPath, _ := cmd.Flags().GetString("destination")
	//accumulative, _ := cmd.Flags().GetBool("accumulative")
	serverFiles, _ := cmd.Flags().GetBool("server-files")
	serverFilesPath, _ := cmd.Flags().GetString("server-files-path")

	absSourcePath, err := filepath.Abs(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	absDestinationPath, err := filepath.Abs(destinationPath)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(absDestinationPath)
	if err == os.ErrNotExist {
		err := os.Mkdir(absDestinationPath, 0700)
		if err != nil {
			log.Fatal(err)
		}
	} else if err == os.ErrExist {
		err := os.RemoveAll(absDestinationPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	ec := ExtractCmd{
		Path:             absDestinationPath,
		StoreServerFiles: serverFiles,
		ServerFilesPath:  serverFilesPath,
	}

	files, err := ioutil.ReadDir(absSourcePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fPath := fmt.Sprintf("%v/%v", absSourcePath, f.Name())
		absFPath, err := filepath.Abs(fPath)
		if err != nil {
			log.Fatal(err)
		}
		ec.extract(absFPath, f.Name())
	}
}

func (ec *ExtractCmd) extract(path string, folderName string) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var gf GDP

	err = restruct.Unpack(data, binary.LittleEndian, &gf)

	if err != nil {
		log.Errorf("%v %v", folderName, err)
	}
	eof := false
	for _, f := range gf.Files {
		if eof {
			return
		}

		var directories string

		if ec.StoreServerFiles {
			directories =  ec.ServerFilesPath + "/" + "/"
		} else {
			directories =  ec.Path + "/" + folderName + "/"
		}

		segments := strings.Split(f.Name, "\\")

		if !ec.StoreServerFiles {
			for i := 1; i < len(segments)-1; i++ {

				if i != 1 {
					directories += "/"
				}

				directories += segments[i]
			}
		}

		isServerFile := make(map[string]bool)
		interestingFiles := "Fiesta.pdb,BH_Helga.lua,BH_Helga.luac,BH_Humar.lua,BH_Humar.luac,Chimera.lua,NPC_AutoLevel.lua,Psy_Psyken.lua,AbStateSaveTypeInfo.shn,AccUpGradeInfo.shn,AccUpgrade.shn,ActionEffectAbState.shn,ActionRangeFactor.shn,ActiveSkillInfoServer.shn,AdminLvSet.shn,ArkAstanica_Skill01_W.bmp,ArkAstron_Skill02_W.bmp,ArkAstron_Skill02_W_Test.bmp,B_Albireo_Skill01_W.bmp,B_Albireo_Skill02_W.bmp,Chimera_Skill_W03_2.BMP,Eglack_Skill05_W.bmp,FinalLight.bmp,KDFargels_DKnight_Skill01_W_G.bmp,KDFargels_DKnight_Skill02_W_G.bmp,KDFargels_Skill01_W_g.bmp,KDFargels_Skill02_W_g.bmp,KDFargels_Spearman_Skill01_W.bmp,LightShot01.bmp,LightShot02.bmp,LightShot03.bmp,SD_DragonSkill05_W.bmp,SD_DragonSkill07_W_g.bmp,SD_DragonSkill07_W_g1.bmp,SD_DragonSkill07_W_g2.bmp,SD_DragonSkill07_W_g3.bmp,SD_DragonSkill07_W_g4.bmp,SD_KingCrabSkill06_W_g.bmp,SD_KingCrabSkill08_W.bmp,SW_FAvanas_Skill01_W.BMP,SW_FAvanas_Skill02_W.BMP,SW_FAvanas_Skill03_W.BMP,SW_FCitrie_Skill07_W.bmp,SW_FFocalor_Skill04_W.BMP,SW_FFocalor_Skill05_N.bmp,SW_FFocalor_Skill05_N_01.bmp,SW_FFocalor_Skill05_N_02.bmp,SW_FFocalor_Skill05_N_03.bmp,SW_FFocalor_Skill05_N_04.bmp,SW_FFocalor_Skill05_N_05.bmp,SW_IFocalor_Skill04_W_01.bmp,SW_IFocalor_Skill04_W_02.bmp,SW_IFocalor_Skill04_W_03.BMP,SW_IFocalor_Skill04_W_04.BMP,SW_IFocalor_Skill04_W_05.BMP,S_Anais_Skill02_W_01.bmp,S_Anais_Skill02_W_02.bmp,S_Anais_Skill02_W_03.bmp,S_FreloanLeg_Skill01_W.bmp,S_Freloan_Skill01_W_01.bmp,S_Freloan_Skill01_W_02.bmp,S_Freloan_Skill01_W_03.bmp,S_Freloan_Skill01_W_04.bmp,S_Freloan_Skill01_W_05.bmp,S_Freloan_Skill02_W.bmp,S_Varamus_Skill02_W.bmp,Thumbs.db,WarH_BossRoom.bmp,AreaSkill.shn,AttendSchedule.shn,BMP.shn,BRAccUpgrade.shn,BelongDice.shn,B_AHarp.shmd,B_Trini.shmd,Eld.shmd,Job2_Dn01.shmd,KDArena.shab,KDCake.shab,KDFargels.shmd,KDHBat1.shab,KDWater.shab,Leviathan.shab,Maghda.shmd,NewTower.shmd,OX_field.shmd,Rou.shmd,RouCos01.shmd,RouCos02.shmd,RouCos03.shmd,RouTemDn01.shab,RouVal01.shmd,Tower02.shab,UrgDragon.shab,UrgDragon.shad,UrgFire01.shmd,CharacterTitleStateServer.shn,CollectCardDropRate.shn,CollectCardMobGroup.shn,CollectCardStarRate.shn,DamageLvGapEVP.shn,DamageLvGapPVE.shn,DamageLvGapPVP.shn,DefaultCharacterData.txt,DiceGame.shn,DiceRate.shn,EnchantSocketRate.shn,FieldLvCondition.shn,FriendPointReward.shn,GBBanTime.shn,GBDiceGame.shn,GBDiceRate.shn,GBEventCode.shn,GBExchangeMaxCoin.shn,GBJoinGameMember.shn,GBReward.shn,GBSMAll.shn,GBSMBetCoin.shn,GBSMCardRate.shn,GBSMCenter.shn,GBSMGroup.shn,GBSMJPRate.shn,GBSMLine.shn,GBSMNPC.shn,GBTaxRate.shn,GTIBreedSubject.shn,GTIGetRate.shn,GTIGetRateGap.shn,GTIServer.shn,GTWinScore.shn,GroupAbState.shn,GuildAcademy.shn,GuildAcademyLevelUp.shn,GuildAcademyRank.shn,GuildGradeData.shn,GuildGradeScoreData.shn,GuildLevelScoreData.shn,GuildTournament.shn,GuildTournamentLvGap.shn,GuildTournamentMasterBuff.shn,GuildTournamentOccupy.shn,GuildTournamentReward.shn,GuildTournamentScore.shn,HolyPromiseReward.shn,ItemDropLog.shn,ItemInfoServer.shn,ItemInvenDel.shn,ItemMerchantInfo.shn,ItemOptions.shn,ItemPackage.shn,ItemServerEquipTypeInfo.shn,ItemShop.shn,ItemSort.shn,ItemUpgrade.shn,ItemUseEffect.shn,JobEquipInfo.shn,KQItem.shn,KingdomQuest.shn.bak,KingdomQuestMap.shn,KingdomQuestRew.shn,LCGroupRate.shn,LCReward.shn,BH_Albireo.lua,BH_Helga.lua,BH_Helga.luac,BH_Humar.lua,BH_Humar.luac,B_Albireo.lua,BallEgg.lua,Chimera.lua,ClassChangeMaster01.lua,ClassChangeMaster02.lua,ClassChangeMaster03.lua,Defense.lua,E_CacaoBud.lua,E_HwinIn.lua,E_HwinOut.lua,E_MomSlime.lua,E_Ski_CongressNPC.lua,E_XTreeBig.lua,Egg2014_GoldEgg.lua,EldCastleLordElderiss.lua,NPC_AutoLevel.lua,Oluming.lua,Psy_Psyken.lua,Toryming.lua,Xiaoming.lua,CrystalCH.lua,Boss.lua,Chat.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,CrystalCastle.lua,Boss.lua,Chat.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,AHarp.lua,AHarp_Data.lua,AdlF.lua,AdlF_Gate.lua,AdlF_Guarder.lua,AdlF_Karen.lua,AdlF_Loussier.lua,AdlF_MagicStone.lua,AdlF_Zone1.lua,AdlF_Zone2.lua,AdlF_Zone3.lua,AdlF_Zone4.lua,AdlFH.lua,AdlFH_Gate.lua,AdlFH_Guarder.lua,AdlFH_Karen.lua,AdlFH_Loussier.lua,AdlFH_MagicStone.lua,AdlFH_Zone1.lua,AdlFH_Zone2.lua,AdlFH_Zone3.lua,AdlFH_Zone4.lua,Bla.lua,Boss.lua,Chat.lua,Name.lua,Process.lua,Regen.lua,Stuff.lua,Progress.lua,Routine.lua,SubFunc.lua,CrystalCH.lua,Boss.lua,Chat.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,CrystalCastle.lua,Boss.lua,Chat.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,Boss.lua,Chat.lua,Name.lua,Process.lua,Regen.lua,Stuff.lua,Progress.lua,Routine.lua,SubFunc.lua,GraveYard.lua,Boss.lua,Chat.lua,Name.lua,Process.lua,Regen.lua,Stuff.lua,Progress.lua,Routine.lua,SubFunc.lua,GraveYardH.lua,Boss.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,IyzelTower.lua,Boss.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,IyzelTowerH.lua,Chat.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,Leviathan.lua,Chat.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,LeviathanH.lua,Chat.lua,Name.lua,Process.lua,Regen.lua,SkillInfo_KingCrab.lua,SkillInfo_KingSlime.lua,SkillInfo_MiniDragon.lua,Progress.lua,Routine.lua,Routine_KingCrab.lua,Routine_KingSlime.lua,Routine_MiniDragon.lua,SubFunc.lua,SD_Vale01.lua,Boss.lua,Chat.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,SecretLab.lua,Boss.lua,Chat.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,SecretLabH.lua,Boss.lua,Chat.lua,Name.lua,Process.lua,Regen.lua,Stuff.lua,Progress.lua,Routine.lua,SubFunc.lua,Siren.lua,Boss.lua,Chat.lua,Name.lua,Process.lua,Regen.lua,Stuff.lua,Progress.lua,Routine.lua,SubFunc.lua,SirenH.lua,WarBL.lua,WarBLData.lua,WarBLDeInitFuntion.lua,WarBLEventMobRoutine.lua,WarBLEventRoutine.lua,WarBLInitFuntion.lua,WarBLH.lua,WarBLHData.lua,WarBLHDeInitFuntion.lua,WarBLHEventMobRoutine.lua,WarBLHEventRoutine.lua,WarBLHInitFuntion.lua,WarH.lua,WarHData.lua,WarHDeInitFunction.lua,WarHEventMobRoutine.lua,WarHEventRoutine.lua,WarHFunction.lua,WarHInitFunction.lua,WarHH.lua,WarHHData.lua,WarHHDeInitFunction.lua,WarHHEventMobRoutine.lua,WarHHEventRoutine.lua,WarHHFunction.lua,WarHHInitFunction.lua,WarL.lua,WarLData.lua,WarLDeInitFuntion.lua,WarLEventMobRoutine.lua,WarLEventRoutine.lua,WarLInitFuntion.lua,WarLH.lua,WarLHData.lua,WarLHDeInitFuntion.lua,WarLHEventMobRoutine.lua,WarLHEventRoutine.lua,WarLHInitFuntion.lua,WarN.lua,WarNData.lua,WarNFunc.lua,WarNRoutine.lua,WarNH.lua,WarNHData.lua,WarNHFunc.lua,WarNHRoutine.lua,AntiHenis.lua,Boss.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,Boss.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,EmperorSlime.lua,Progress.lua,Routine.lua,SubFunc.lua,Boss.lua,ItemDrop.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,GoldHill.lua,Boss.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,HMiniDragon.lua,Monster.lua,NPC.lua,Process.lua,Name.lua,Regen.lua,Name.lua,Regen.lua,Name.lua,Regen.lua,Name.lua,Regen.lua,Name.lua,Regen.lua,Name.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,KDArena.lua,KDArena1.lua,KDArena2.lua,KDArena3.lua,KDArena4.lua,KDArena5.lua,KDArena6.lua,Name.lua,Process.lua,Regen.lua,Servant.lua,Progress.lua,Routine.lua,SubFunc.lua,KDCake.lua,KDEgg.lua,KDEggData.lua,KDEggFunc.lua,KDEggObjectRoutine.lua,Boss.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,KDFargels.lua,KDMine.lua,KDMineData.lua,KDMineFunc.lua,KDMineObjectRoutine.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,KDSoccer.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,KDSoccer_W.lua,KDSpring.lua,KDSpring_Data.lua,KDSpring_StepFunc.lua,Name.lua,Process.lua,Regen.lua,Servant.lua,Progress.lua,Routine.lua,SubFunc.lua,KDWater.lua,Boss.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,KingSlime.lua,Boss.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,KingSlime2.lua,Boss.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,KingSlime3.lua,Boss.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,KingSlime4.lua,Boss.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,Kingkong.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,LegendOfBijou.lua,Boss.lua,NPC.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,MaraPirate.lua,Boss.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,MiniDragon.lua,Name.lua,Process.lua,Regen.lua,Dice.lua,Progress.lua,Routine.lua,SubFunc.lua,Infect.lua,Infect_Data.lua,Infect.lua,Infect_Data.lua,MegaMob.lua,MegaMob_Data.lua,MegaMob.lua,MegaMob_Data.lua,MegaMob.lua,MegaMob_Data.lua,RouN.lua,PetBase.lua,PetBaseActionData.lua,PetBaseIdleActionFunc.lua,PetBaseRoutineFunc.lua,PetBaseSettingFunc.lua,PetBaseSubFunc.lua,PetCommon.lua,SubFunc.lua,Chat.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,Job2_Forest.lua,Name.lua,Process.lua,Regen.lua,Progress.lua,Routine.lua,SubFunc.lua,Job2_Gamb.lua,Tutorial.lua,TutorialData.lua,common.lua,MapBuff.shn,MiniHouseDummy.shn,MobAbStateDropSetting.shn,AdlFH_Eglack.txt,AdlFH_EglackMad.txt,AdlFH_Fknuckleman.txt,AdlFH_Salare.txt,AdlF_Fknuckleman.txt,AngelicKaren.txt,Anti_Henis_A100.txt,Anti_Henis_A101.txt,Anti_Henis_A102.txt,Anti_Henis_A103.txt,Anti_Henis_A104.txt,Anti_Henis_A105.txt,Anti_Henis_A106.txt,Anti_Henis_A107.txt,Anti_Henis_A108.txt,Anti_Henis_A109.txt,Anti_Henis_A110.txt,Anti_Henis_A60.txt,Anti_Henis_A61.txt,Anti_Henis_A62.txt,Anti_Henis_A63.txt,Anti_Henis_A64.txt,Anti_Henis_A65.txt,Anti_Henis_A66.txt,Anti_Henis_A67.txt,Anti_Henis_A68.txt,Anti_Henis_A69.txt,Anti_Henis_A70.txt,Anti_Henis_A90.txt,Anti_Henis_A92.txt,Anti_Henis_A94.txt,Anti_Henis_A95.txt,Anti_Henis_A96.txt,Anti_Henis_A98.txt,Anti_Henis_A99.txt,Anti_Henis_C100.txt,Anti_Henis_C101.txt,Anti_Henis_C102.txt,Anti_Henis_C103.txt,Anti_Henis_C104.txt,Anti_Henis_C105.txt,Anti_Henis_C106.txt,Anti_Henis_C107.txt,Anti_Henis_C108.txt,Anti_Henis_C109.txt,Anti_Henis_C110.txt,Anti_Henis_C60.txt,Anti_Henis_C61.txt,Anti_Henis_C62.txt,Anti_Henis_C63.txt,Anti_Henis_C64.txt,Anti_Henis_C65.txt,Anti_Henis_C66.txt,Anti_Henis_C67.txt,Anti_Henis_C68.txt,Anti_Henis_C69.txt,Anti_Henis_C70.txt,Anti_Henis_C90.txt,Anti_Henis_C94.txt,Anti_Henis_C95.txt,Anti_Henis_C96.txt,Anti_Henis_C97.txt,Anti_Henis_C98.txt,Anti_Henis_C99.txt,Anti_Henis_F100.txt,Anti_Henis_F101.txt,Anti_Henis_F102.txt,Anti_Henis_F103.txt,Anti_Henis_F104.txt,Anti_Henis_F105.txt,Anti_Henis_F106.txt,Anti_Henis_F107.txt,Anti_Henis_F108.txt,Anti_Henis_F109.txt,Anti_Henis_F110.txt,Anti_Henis_F60.txt,Anti_Henis_F61.txt,Anti_Henis_F62.txt,Anti_Henis_F63.txt,Anti_Henis_F64.txt,Anti_Henis_F65.txt,Anti_Henis_F66.txt,Anti_Henis_F67.txt,Anti_Henis_F68.txt,Anti_Henis_F69.txt,Anti_Henis_F70.txt,Anti_Henis_F90.txt,Anti_Henis_F92.txt,Anti_Henis_F94.txt,Anti_Henis_F95.txt,Anti_Henis_F96.txt,Anti_Henis_F97.txt,Anti_Henis_F98.txt,Anti_Henis_F99.txt,Anti_Henis_G_A10.txt,Anti_Henis_G_A100.txt,Anti_Henis_G_A110.txt,Anti_Henis_G_A120.txt,Anti_Henis_G_A20.txt,Anti_Henis_G_A30.txt,Anti_Henis_G_A40.txt,Anti_Henis_G_A50.txt,Anti_Henis_G_A60.txt,Anti_Henis_G_A70.txt,Anti_Henis_G_A80.txt,Anti_Henis_G_A90.txt,Anti_Henis_G_C10.txt,Anti_Henis_G_C100.txt,Anti_Henis_G_C110.txt,Anti_Henis_G_C120.txt,Anti_Henis_G_C20.txt,Anti_Henis_G_C30.txt,Anti_Henis_G_C40.txt,Anti_Henis_G_C50.txt,Anti_Henis_G_C60.txt,Anti_Henis_G_C70.txt,Anti_Henis_G_C80.txt,Anti_Henis_G_C90.txt,Anti_Henis_G_F10.txt,Anti_Henis_G_F100.txt,Anti_Henis_G_F110.txt,Anti_Henis_G_F120.txt,Anti_Henis_G_F20.txt,Anti_Henis_G_F30.txt,Anti_Henis_G_F40.txt,Anti_Henis_G_F50.txt,Anti_Henis_G_F60.txt,Anti_Henis_G_F70.txt,Anti_Henis_G_F80.txt,Anti_Henis_G_F90.txt,Anti_Henis_G_M10.txt,Anti_Henis_G_M100.txt,Anti_Henis_G_M110.txt,Anti_Henis_G_M120.txt,Anti_Henis_G_M20.txt,Anti_Henis_G_M30.txt,Anti_Henis_G_M40.txt,Anti_Henis_G_M50.txt,Anti_Henis_G_M60.txt,Anti_Henis_G_M70.txt,Anti_Henis_G_M80.txt,Anti_Henis_G_M90.txt,Anti_Henis_M100.txt,Anti_Henis_M101.txt,Anti_Henis_M102.txt,Anti_Henis_M103.txt,Anti_Henis_M104.txt,Anti_Henis_M105.txt,Anti_Henis_M106.txt,Anti_Henis_M107.txt,Anti_Henis_M108.txt,Anti_Henis_M109.txt,Anti_Henis_M110.txt,Anti_Henis_M60.txt,Anti_Henis_M61.txt,Anti_Henis_M62.txt,Anti_Henis_M63.txt,Anti_Henis_M64.txt,Anti_Henis_M65.txt,Anti_Henis_M66.txt,Anti_Henis_M67.txt,Anti_Henis_M68.txt,Anti_Henis_M69.txt,Anti_Henis_M70.txt,Anti_Henis_M90.txt,Anti_Henis_M91.txt,Anti_Henis_M92.txt,Anti_Henis_M95.txt,Anti_Henis_M96.txt,Anti_Henis_M97.txt,Anti_Henis_M98.txt,Anti_Henis_M99.txt,BH_Albireo.txt,BH_Guardian.txt,BH_Helga.txt,BH_Humar.txt,BH_Looter.txt,B_Albireo.txt,B_CrackerGuardian.txt,B_CrackerHumar.txt,B_CrackerLooter.txt,BomBoogy01.txt,BomBoogy02.txt,BomBoogy03.txt,BomBoogy04.txt,BoogyGuardian.txt,C_JewelGolem.txt,Chimera.txt,DT_FFocalor.txt,DT_FFocalor_C.txt,DT_IFocalor_C.txt,DT_Ifocalor.txt,DT_SFocalor.txt,DT_SFocalor_C.txt,DT_TFocalor.txt,DT_TFocalor_C.txt,Eglack.txt,EglackMad.txt,EmperorCrab.txt,Event_H_MiniDragon.txt,FireTotem.txt,Firepamelia.txt,Helga.txt,ID_BigMudMan.txt,ID_EarthCalerben.txt,ID_EarthNerpa.txt,ID_FandomCornelius.txt,ID_FireShella.txt,ID_FireTaitan.txt,ID_FlameSpirit.txt,ID_GiantMagmaton.txt,ID_Kruge.txt,ID_NestAlca.txt,ID_NestBaridon.txt,ID_NestGuardian.txt,ID_NestMadSlug.txt,ID_NestWeasel.txt,ID_Weasel.txt,KDFargels.txt,KQ_H_MiniDragon.txt,KQ_KalBanObeb.txt,Karen.txt,KillerHide.txt,KingBoogy.txt,LabH_19.txt,LabH_20.txt,LabH_23.txt,LabH_25.txt,Lab_19.txt,Lab_20.txt,Lab_23.txt,Lab_25.txt,LevH_EmperorCrab.txt,LevH_ID_NestAlca.txt,LevH_ID_NestBaridon.txt,LevH_KingBoogy.txt,LevH_ViciousLeviathan.txt,LevH_ViciousLeviathan01.txt,Maghda.txt,MasicStaff.txt,NT_C_JewelGolem.txt,NT_DT_IFocalor.txt,NT_DT_SFocalor.txt,NT_DT_TFocalor.txt,NT_Eglack.txt,NT_EmperorCrab.txt,NT_ID_BigMudMan.txt,NT_ID_FandomCornelius.txt,NT_ID_FireTaitan.txt,NT_ID_GiantMagmaton.txt,NT_ID_Weasel.txt,NT_KingBoogy.txt,NT_Lab_19.txt,NT_Lab_20.txt,NT_Lab_23.txt,NT_Lab_25.txt,NT_Psyken.txt,NT_T_DustGolem.txt,NT_T_IronGolem.txt,NT_T_PoisonGolem.txt,NT_T_StoneGolem.txt,NT_ViciousLeviathan01.txt,P_Psy_Mist1.txt,P_Psy_Mist2.txt,P_Psy_Mist3.txt,Psy_Mist.txt,Psy_Psyken.txt,Psy_Wraith.txt,Psyken.txt,SD_Dragon.txt,SD_KingCrab.txt,SD_KingSlime.txt,Salare.txt,Silberk.txt,T_ArchMageBook00.txt,T_ArchMageBook01.txt,T_Boar.txt,T_DesertWolf.txt,T_DustGolem.txt,T_FlyingStaff00.txt,T_FlyingStaff01.txt,T_GangImp.txt,T_Ghost.txt,T_HungryWolf.txt,T_IceViVi.txt,T_Imp.txt,T_IronGolem.txt,T_IronSlime00.txt,T_IronSlime01.txt,T_Kamaris01.txt,T_Kamaris02.txt,T_Kebing.txt,T_KingCall.txt,T_KingSpider.txt,T_OldFox.txt,T_PoisonGolem.txt,T_Prock.txt,T_Ratman.txt,T_SkelArcher00.txt,T_SkelArcher01.txt,T_SkelArcher02.txt,T_SkelWarrior.txt,T_Skeleton.txt,T_Spider00.txt,T_Spider01.txt,T_StoneGolem.txt,T_Zombie.txt,TestFireTotem.txt,TestSilberk02.txt,UrgDTH_ID_BigMudMan.txt,UrgDTH_ID_EarthCalerben.txt,UrgDTH_ID_EarthNerpa.txt,UrgDTH_ID_FandomCornelius.txt,UrgDTH_ID_FireShella.txt,UrgDTH_ID_FireTaitan.txt,UrgDTH_ID_FlameSpirit.txt,UrgDTH_ID_GiantMagmaton.txt,UrgDTH_ID_Kruge.txt,UrgDTH_ID_Weasel.txt,ViciousLeviathan.txt,ViciousLeviathan01.txt,WarH_FAvanas2.txt,MobAutoAction.shn,DefaultBehavior.ps,BossRobo.ps,EndlessMaze.ps,GordonMaster.ps,KingSlime.ps,MaraPirate.ps,TravelerDungeon.ps,UnderHall.ps,MobConditionServer.shn,MobInfoServer.shn,MobKillAble.shn,MobKillAnnounce.shn,MobKillLog.shn,MobLifeTime.shn,Adl.txt,AdlThorn01.txt,AdlVal01.txt,AlDn01.txt,AlDn02.txt,ArkDn01.txt,ArkDn02.txt,BH_Albi.txt,BH_Cracker.txt,BH_Helga.txt,B_Albi.txt,B_Cracker.txt,BerFrz01.txt,BerKal01.txt,BerVale01.txt,BerVale02.txt,Bera.txt,CemDn01.txt,CemDn02.txt,E_Hwin.txt,E_Olympic.txt,EchoCave.txt,Eld.txt,EldCem01.txt,EldCem02.txt,EldFor01.txt,EldGbl01.txt,EldGbl02.txt,EldPri01.txt,EldPri02.txt,EldPriDn01.txt,EldPriDn02.txt,EldSleep01.txt,ElfDn01.txt,ElfDn02.txt,Fbattle01.txt,Fbattle02.txt,Fbattle03.txt,FireDn01.txt,FireDn02.txt,ForDn01.txt,ForDn02.txt,GBHouse.txt,GblDn01.txt,GblDn02.txt,GoldCave.txt,GuildT0400.txt,GuildT0401.txt,GuildT0402.txt,GuildT0403.txt,GuildT0404.txt,GuildT0405.txt,GuildT0406.txt,GuildT0407.txt,AdlF.txt,AdlFH.txt,Leviathan.txt,NewTower.txt,Siren.txt,Tower01.txt,Tower02.txt,Tower03.txt,UrgDragon.txt,WarN.txt,KDHBat1.txt,Job2_Dn01.txt,Job2_Dn02.txt,KDAntiHenis.txt,KDEchoCave.txt,KDEddyHill(old).txt,KDEddyHill.txt,KDEddyHill2.txt,KDEddyHill3.txt,KDEddyHill4.txt,KDEnMaze.txt,KDFargels.txt,KDGoldHill(Making).txt,KDGoldHill.txt,KDGreenHill.txt,KDHBat1.txt,KDHDragon.txt,KDHero.txt,KDHoneying.txt,KDKingkong.txt,KDKingkong2.txt,KDKingkong3.txt,KDMDragon.txt,KDPrtShip.txt,KDRockCan.txt,KDSoccer_W.txt,KDSpider.txt,KDTrDn.txt,KDUnHall.txt,KDUnHall2.txt,KDVictor.txt,KD_Kingkong.txt,KQ_HONEYING.txt,SlimeKQHig.txt,SlimeKQLow.txt,SlimeKQMed.txt,KingdomQuest.txt,Linkfield01.txt,Linkfield02.txt,PriDn01.txt,PriDn02.txt,PsyIn.txt,PsyInDn02.txt,PsyOut.txt,QField01.txt,QField02.txt,QField03.txt,QField04.txt,R_Helga01.txt,Rou.txt,RouCos01.txt,RouCos02.txt,RouCos03.txt,RouN.txt,RouTemDn01.txt,RouTemDn02.txt,RouVal01.txt,RouVal02.txt,SwaDn01.txt,SwaDn02.txt,Urg.txt,UrgDark01.txt,UrgFire01.txt,UrgFireDn01.txt,UrgSwa01.txt,UrgSwaDn01.txt,Urg_Alruin.txt,ValDn01.txt,ValDn02.txt,WindyCave.txt,desktop.ini,MobRegenAni.shn,MobResist.shn,AdlThornR01.txt,AdlThornR02.txt,AdlThornR03.txt,AdlThornR04.txt,AdlThornR05.txt,AdlThornR06.txt,AdlThornR07.txt,AdlThornR08.txt,AdlThornR09.txt,AdlThornR10.txt,AdlThornR11.txt,AdlThornR12.txt,BerValeDw01.txt,BerValeDw02.txt,BerValeDw03.txt,BerValeDw04.txt,BerValeGaruda.txt,GB_Waitress01.txt,GB_Waitress02.txt,GB_Waitress03.txt,GB_Waitress04.txt,GB_Waitress05.txt,GB_Waitress06.txt,GB_Waitress07.txt,GB_Waitress08.txt,GB_Waitress09.txt,GB_Waitress10.txt,Kal01_HERO1.txt,Action.xls,B_Slime.txt,BallCake01.txt,BallWater.txt,C_Gate01.txt,ChristmasTree.txt,DT_EntranceGate.txt,DT_ExitGate.txt,DT_RadionOre.txt,E_SixYear_Dance.txt,E_SkiFlag_Blue.txt,E_SkiFlag_Gold.txt,E_SkiFlag_Red.txt,E_Ski_IDHoneying.txt,E_Ski_Snowman.txt,Egg2014_BigEgg.txt,GTI_BoxAll.txt,GTI_BoxTeamA.txt,GTI_BoxTeamB.txt,Gate_AdlF.txt,Gate_Lab.txt,IDLeviathanGate01.txt,IDMapLinkGate00.txt,IDMapLinkGate01.txt,IDMapLinkGate02.txt,IDMapLinkGate03.txt,KDSoccer_Ball.txt,KDSoccer_Ball_14.txt,KDSoccer_Invincible.txt,KDSoccer_SpeedUp.txt,KarenGate.txt,LightField01.txt,LightField02.txt,LightField03.txt,LightField04.txt,LightField05.txt,LightOrb01.txt,LightOrb02.txt,LightOrb03.txt,LightOrb04.txt,LightOrb05.txt,LightOrb06.txt,LightOrb07.txt,LightOrb08.txt,LightOrb09.txt,LightOrb10.txt,MapLinkGate.txt,MapLinkGate01.txt,MultiProtect.txt,MultiProtect02.txt,MultiProtect03.txt,SpImShield.txt,SpUpShoes.txt,T_Gate01.txt,T_Gate02.txt,WarBL_EntranceGate.txt,WarH_EntranceGate.txt,MobSpecies.shn,MobWeapon.shn,MsgWorldManager.shn,MultiHitType.shn,MysteryVaultServer.shn,AdlAertsina.txt,AdlLoussier.txt,AdlSkillEdwina.txt,AdlSmithAlexia.txt,AlruinItemMctGeric.txt,AlruinSkillPaela.txt,AlruinSmithMacurdos.txt,BeraGuardArcher.txt,BeraItemEdmong.txt,BeraItemMilly.txt,BeraSkillHal.txt,BeraSmithMcDilan.txt,CardRefunder.txt,Card_Mct_1Piece.txt,Card_Mct_3Piece.txt,Card_Mct_BackTail.txt,Card_Mct_FaceHatMask.txt,Card_Mct_Pet.txt,Chaoming.txt,CustomProdNPC.txt,DailyCoinMerchant.txt,Daliy_Merchant.txt,E_HwinQuest.txt,E_Ski_MerchantNPC.txt,E_Ski_RentMachine.txt,E_XXiaoming.txt,Egg2014_HoshemingNPC.txt,Egg_Digger.txt,EldArcGuard03.txt,EldFurnitureForestTall.txt,EldItemMctKenton.txt,EldItemMctNina.txt,EldItemMctNina2.txt,EldPalSkillKeest.txt,EldScoSkillDeikid.txt,EldSmithKarls.txt,EldWarSkillMarty.txt,EldWeaponTitleMctBran.txt,EldWizSkillWishis.txt,GB_CoinStore.txt,GB_MasterRoan.txt,GuildItemMct.txt,HednisSkillGrunt.txt,HednisSmithRohan.txt,IDVaultNPC.txt,IM_Arena01.txt,IM_Arena02.txt,IM_Arena_TE.txt,ItemMctJelluin.txt,Joker.txt,KDSoccer_MctNPC.txt,KDSoccer_MctNPC_14.txt,KQSpring_Bman.txt,KQSpring_Rman.txt,KarenMct.txt,LC_Machine.txt,LC_MachineBlue.txt,LC_MachineRed.txt,Mct_1Piece.txt,Mct_3Piece.txt,Mct_Back.txt,Mct_Face.txt,Mct_Hat.txt,Mct_House.txt,Mct_Mount.txt,Mct_Pet.txt,Mct_Pet2.txt,Mct_Skin.txt,Mct_Tail.txt,MineDigger.txt,NPCItemList.txt,RouFurnitureForestTom.txt,RouItemMctPey.txt,RouSkillRubi.txt,RouSmithJames.txt,RouSoulMctJulia.txt,RouT_Skill.txt,RouT_Smith.txt,RouTownChiefRoumenus.txt,RouWeaponTitleMctZach.txt,SD_Futureming.txt,Swimming.txt,SwimmingB.txt,SwimmingR.txt,TempSkill.txt,Tiros.txt,UruFurnitureForestTeem.txt,UruItemMctVellon.txt,UruSkillChyburn.txt,UruSmithHans.txt,WeddingDreian.txt,XiaomingB_7th.txt,XiaomingR_7th.txt,Xiaoming_7th.txt,PSkillSetAbstate.shn,PartyBonusByLvDiff.shn,PartyBonusByMember.shn,PartyBonusLimit.shn,PupCase.shn,PupCaseDesc.shn,PupFactorCondition.shn,PupMind.shn,PupPriority.shn,PupServer.shn,QuestScript.shn,QuestSpecies.shn,RandomOption.shn,RandomOptionCount.shn,RareMoverEachRate.shn,RareMoverRate.shn,RareMoverSubRate.shn,ReactionType.shn,GuildTournament.ps,GuildTournament1.ps,Maghda.ps,NewTower.ps,d_NestOfLeviathan.ps,d_graveyard.ps,ConditionOfHero.ps,ConditionOfHero.ps,GordonMaster.ps,HMiniDragon.ps,Honeying.ps,KQHBat1.ps,KQHBat2.ps,KQHBat3.ps,KQHBat4.ps,KQHBat5.ps,KickOut.ps,Fail.ps,Main.ps,Suc.ps,Kingkong.ps,MiniDragon.ps,Quest.ps,RoumenGate.ps,TesScript.ps,Fail.ps,Main.ps,Suc.ps,UnderHall(Bak).ps,UnderHall.ps,UnderHall2.ps,AdlThorn01.ps,AlDn01.ps,BerFrz01.ps,CemDn01.ps,EldCem01.ps,EldCem02.ps,EldFor01.ps,EldGbl01.ps,EldPriDn02.ps,EldSleep01.ps,Elderine.ps,ElfDn01.ps,FireDn01.ps,ForDn01.ps,GblDn01.ps,Monkey.ps,MonkeyBreed01.ps,MonkeyBreed02.ps,MonkeyBreed03.ps,OXFieldInit.ps,PriDn01.ps,QField03.ps,QField04.ps,RouCos03.ps,RouTemDn02.ps,RouVal02.ps,Roumen.ps,UrgFire01.ps,UrgFireDn01.ps,UrgSwa01.ps,UrgSwaDn01.ps,ValDn01.ps,Job2_Gamb.ps,JobChange1.ps,JobChange2-1.ps,ConditionOfHero.ps,GoldHill.ps,KickOut.ps,Kingkong.ps,MaraPirate.ps,MiniDragon.ps,Quest.ps,TesScript.ps,UnderHall.ps,Wedding.ps,AdlF.txt,AdlFH.txt,D_Graveyard.txt,Defense.txt,Defense01.txt,ETC.txt,EldEvent.txt,Event.txt,GraveYard.txt,GraveYardH.txt,Honeying.txt,JC100.txt,JobChange1.txt,JobChange2-1.txt,JobChange2-2.txt,JobChange2-3.txt,KDArena.txt,KDEgg.txt,KDFargels.txt,KDGreenHill.txt,KDMine.txt,KDSpring.txt,KQAntiHenis100.txt,KQConditionOfHero.txt,KQGoldHill.txt,KQGordonMaster.txt,KQHBat1.txt,KQHBat2.txt,KQHBat3.txt,KQHBat4.txt,KQHBat5.txt,KQHoneying.txt,KQKingSlime.txt,KQKingkong.txt,KQMaraPirate.txt,KQMiniDragon.txt,KQUnderHall.txt,KQUnderHall2.txt,LegendOfBijou.txt,MapName.txt,MenuString.txt,Scenario.txt,Script.txt,Siren.txt,Tower01.txt,Tower02.txt,Tower03.txt,WarBL.txt,WarH.txt,WarL.txt,WarN.txt,Wedding.txt,d_NestOfLeviathan.txt,SetItem.shn,SetItemEffect.shn,ShineReward.shn,SpamerPenalty.shn,SpamerPenaltyRule.shn,SpamerReport.shn,StateField.shn,StateItem.shn,StateMob.shn,ToggleSkill.shn,TutorialCharacterData.txt,SetItemView.shn,ChrCommon.txt,DamageByAngle.txt,DamageBySoul.txt,ExpRecalculation.txt,Field.txt,ItemDropGroup.txt,ItemDropTable.txt,ItemOptions.txt,ItemUseFunction.txt,MiscDataTable.txt,MobChat.txt,NPC.txt,NPCAction.txt,ParamArcherServer.txt,ParamAssassinServer.txt,ParamChaserServer.txt,ParamClericServer.txt,ParamCleverFighterServer.txt,ParamCloserServer.txt,ParamCruelServer.txt,ParamEnchanterServer.txt,ParamFighterServer.txt,ParamGladiatorServer.txt,ParamGuardianServer.txt,ParamHawkArcherServer.txt,ParamHighClericServer.txt,ParamHolyKnightServer.txt,ParamJokerServer.txt,ParamKnightServer.txt,ParamMageServer.txt,ParamPaladinServer.txt,ParamRangerServer.txt,ParamSaviorServer.txt,ParamScoutServer.txt,ParamSentinelServer.txt,ParamSharpShooterServer.txt,ParamWarriorServer.txt,ParamWarrockServer.txt,ParamWizMageServer.txt,ParamWizardServer.txt,PineScript.txt,Quest.txt,QuestParser.txt,RandomOptionTable.txt,RecallCoord.txt,SubLayerInteract.txt,TreasureReward.txt,"
		for _, f := range strings.Split(interestingFiles, ",") {
			isServerFile[f] = true
		}

		fileName := "/"
		if ec.StoreServerFiles {
			if _, ok := isServerFile[segments[len(segments)-1]]; ok {
				fileName += gf.Version + "_" + ksuid.New().String() + "_" + segments[len(segments)-1]
			} else {
				return
			}
		} else {
			if segments[len(segments)-1] == "Fiesta.bin" {
				fileName += gf.Version + "_" + ksuid.New().String() + "_" + segments[len(segments)-1]
			} else {
				fileName += segments[len(segments)-1]
			}
		}

		absPath, err := filepath.Abs(directories)
		if err != nil {
			log.Error(err)
			return
		}

		err = os.MkdirAll(absPath, 0700)

		if err != nil {
			log.Error(err)
			return
		}

		fileAbsPath, err := filepath.Abs(directories+fileName)
		if err != nil {
			log.Error(err)
			return
		}

		_, err = os.Stat(fileAbsPath)

		if err == nil {
			log.Infof("file %v already exists, skipping.", fileAbsPath)
			return
		}

		file, err := os.OpenFile(fileAbsPath, os.O_RDONLY|os.O_CREATE, 0700)

		if err != nil {
			fmt.Println(err)
		}

		var b []byte
		if f.Offset+f.Size > uint64(len(data)) {
			b = append(b, data[f.Offset:]...)
			// files are listed but there is no data to read from
			// this means the full patch is the next gdp
			// no clue why this was made like this x.x
			eof = true
		} else {
			b = append(b, data[f.Offset:f.Offset+f.Size]...)
		}

		_, err = file.Write(b)
		if err != nil {
			fmt.Println(err)
		}
	}
}