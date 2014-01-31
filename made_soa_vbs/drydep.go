package made_soa_vbs

import "math"

// Calculates dry deposition velocity, adapted from WRF/Chem 
// module_aerosols_soa_vbs.F subroutine VDVG_2. Original description follows.
// 
//---------------------------------------------------------------------------
//
// *** this routine calculates the dry deposition and sedimentation
//     velocities for the three modes. 
//   Stu McKeen 10/13/08
//   Gaussian Quadrature numerical integration over diameter range for each mode.
// Quadrature taken from Abramowitz and Stegun (1974), equation 25.4.46 and Table 25.10
// Quadrature points are the zeros of Hermite polynomials of order NGAUSdv
//   Numerical Integration allows more complete discription of the
//   Cunningham Slip correction factor, Interception Term (not included previously),
//   and the correction due to rebound for higher diameter particles.
//   Sedimentation velocities the same as original Binkowski code, also the
//   Schmidt number and Brownian diffusion efficiency dependence on Schmidt number the
//   same as Binkowski.
//   Stokes number, and efficiency dependence on Stokes number now according to
//   Peters and Eiden (1992).  Interception term taken from Slinn (1982) with
//   efficiency at .2 micron diam. (0.3%) tuned to yield .2 cm/s deposition velocitiy
//   for needleaf evergreen trees (Pryor et al., Tellus, 2008). Rebound correction
//   term is that of Slinn (1982)
//
//     Original code 1/23/97 by Dr. Francis S. Binkowski. Follows 
//     FSB's original method, i.e. uses Jon Pleim's expression for deposition
//     velocity but includes Marv Wesely's wstar contribution. 
//ia eliminated Stokes term for coarse mode deposition calcs.,
//ia see comments below
// 
// Input variables: Air temperature (BLKTA; [K]), 
func drydep( BLKTA, BLKDENS,                           &
             RA, USTAR, PBLH, ZNTT, RMOLM,  AMU,       &
             DGNUC, DGACC, DGCOR, XLM,                 &
             KNNUC, KNACC,KNCOR,                       &
             PDENSN, PDENSA, PDENSC,                   &
             VSED, VDEP)

// *** calculate size-averaged particle dry deposition and 
//     size-averaged sedimentation velocities.
//     IMPLICIT NONE

      REAL BLKTA( BLKSIZE )         // Air temperature [ K ]
      REAL BLKDENS(BLKSIZE)         // Air density  [ kg m^-3 ]
      REAL RA(BLKSIZE )             // aerodynamic resistance [ s m**-1 ]
      REAL USTAR( BLKSIZE )         // surface friction velocity [ m s**-1 ]
      REAL PBLH( BLKSIZE )          // PBL height (m)
      REAL ZNTT( BLKSIZE )          // Surface roughness length (m)
      REAL RMOLM( BLKSIZE )         // Inverse of Monin-Obukhov length (1/m)
      REAL AMU( BLKSIZE )           // atmospheric dynamic viscosity [ kg m**-1 s**-1 ]
      REAL XLM( BLKSIZE )           // mean free path of dry air [ m ]
      REAL DGNUC( BLKSIZE )         // nuclei mode mean diameter [ m ]
      REAL DGACC( BLKSIZE )         // accumulation  
      REAL DGCOR( BLKSIZE )         // coarse mode
      REAL KNNUC( BLKSIZE )         // nuclei mode Knudsen number 
      REAL KNACC( BLKSIZE )         // accumulation  
      REAL KNCOR( BLKSIZE )         // coarse mode
      REAL PDENSN( BLKSIZE )        // average particle density in nuclei mode [ kg / m**3 ]
      REAL PDENSA( BLKSIZE )        // average particle density in accumulation mode [ kg / m**3 ]
      REAL PDENSC( BLKSIZE )        // average particle density in coarse mode [ kg / m**3 ]

// *** deposition and sedimentation velocities

      REAL VDEP( BLKSIZE, NASPCSDEP) // sedimentation velocity [ m s**-1 ]
      REAL VSED( BLKSIZE, NASPCSSED) // deposition  velocity [ m s**-1 ]

      INTEGER LCELL,N
      REAL DCONST1, DCONST2, DCONST3, DCONST3N, DCONST3A,DCONST3C
      REAL UTSCALE,CZH   // scratch functions of USTAR and WSTAR.
      REAL NU            //kinematic viscosity [ m**2 s**-1 ]
      REAL BHAT
      PARAMETER( BHAT =  1.246 ) // Constant from Binkowski-Shankar approx to Cunningham slip correction.
      REAL COLCTR_BIGD,COLCTR_SMALD
      PARAMETER ( COLCTR_BIGD=2.E-3,COLCTR_SMALD=20.E-6 )  // Collector diameters in Stokes number and Interception Efficiency (Needleleaf Forest)
      REAL SUM0, SUM3, DQ, KNQ, CUNQ, VSEDQ, SCQ, STQ, RSURFQ, vdplim
      REAL Eff_dif, Eff_imp, Eff_int, RBcor
      INTEGER ISTOPvd0,IdoWesCor
      PARAMETER (ISTOPvd0 = 0)  // ISTOPvd0 = 1 means dont call VDVG, particle dep. velocities are set = 0; ISTOPvd0 = 0 means do depvel calcs.

      // no Wesley deposition, otherwise EC is too low
      PARAMETER (IdoWesCor = 0) // IdoWesCor = 1 means do Wesley (85) convective correction to PM dry dep velocities; 0 means don't do correction
      IF (ISTOPvd0.EQ.1)THEN
      RETURN
      ENDIF
// *** check layer value. 

      IF(iprnt.eq.1) print *,'In VDVG, Layer=',LAYER
         IF ( LAYER .EQ. 1 ) THEN // calculate diffusitities and sedimentation velocities
                 
         DO LCELL = 1, NUMCELLS
            DCONST1 = BOLTZ * BLKTA(LCELL) /                                         &
                    ( THREEPI * AMU(LCELL) )
            DCONST2 = GRAV / ( 18.0 * AMU(LCELL) )
            DCONST3 =  USTAR(LCELL)/(9.*AMU(LCELL)*COLCTR_BIGD)
 
// *** now calculate the deposition velocities at layer 1

         NU = AMU(LCELL) / BLKDENS(LCELL) 

         UTSCALE =  1.
        IF (IdoWesCor.EQ.1)THEN
// Wesley (1985) Monin-Obukov dependence for convective conditions (SAM 10/08)
           IF(RMOLM(LCELL).LT.0.)THEN
                CZH = -1.*PBLH(LCELL)*RMOLM(LCELL)
                IF(CZH.GT.30.0)THEN
                  UTSCALE=0.45*CZH**0.6667
                ELSE
                  UTSCALE=1.+(-300.*RMOLM(LCELL))**0.6667
                ENDIF
           ENDIF
        ENDIF   // end of (IdoWesCor.EQ.1) test

        UTSCALE = USTAR(LCELL)*UTSCALE
      IF(iprnt.eq.1)THEN
          print *,'NGAUSdv,xxlsga,USTAR,UTSCALE'
          print *,NGAUSdv,xxlsga,USTAR(LCELL),UTSCALE
          print *,'DCONST2,PDENSA,DGACC,GRAV,AMU'
          print *,DCONST2,PDENSA(LCELL),DGACC(LCELL),GRAV,AMU(LCELL)
      endif
      
// *** nuclei mode 
      
        SUM0=0.
        SUM3=0.
        DO N=1,NGAUSdv
         DQ=DGNUC(LCELL)*EXP(Y_GQ(N)*sqrt2*xxlsgn)  // Diameter (m) at quadrature point
            KNQ=2.*XLM(LCELL)/DQ  // Knudsen number at quadrature point
            CUNQ=1.+KNQ*(1.257+.4*exp(-1.1/KNQ))  // Cunningham correction factor; Pruppacher and Klett (1980) Eq (12-16)
            VSEDQ=PDENSN(LCELL)*DCONST2*CUNQ*DQ*DQ  // Gravitational sedimentation velocity m/s
            SCQ=NU*DQ/DCONST1/CUNQ  // Schmidt number, Brownian diffusion parameter - Same as Binkowski and Shankar
            Eff_dif=SCQ**(-TWO3)    // Efficiency term for diffusion - Same as Binkowski and Shankar
            STQ=DCONST3*PDENSN(LCELL)*DQ**2  // Stokes number, Peters and Eiden (1992)
            Eff_imp=(STQ/(0.8+STQ))**2   // Efficiency term for impaction - Peters and Eiden (1992)
    //       Eff_int=0.3*DQ/(COLCTR_SMALD+DQ) // Slinn (1982) Interception term, 0.3 prefac insures .2 cm/s at .2 micron diam.
            Eff_int=(0.00116+0.0061*ZNTT(LCELL))*DQ/1.414E-7 // McKeen(2008) Intercptn trm, val of .00421 @ ustr=0.475, diam=.1414 micrn, stable, needleleaf evergreen
            RBcor=1. // Rebound correction factor
            vdplim=UTSCALE*(Eff_dif+Eff_imp+Eff_int)*RBcor
    //       vdplim=.002*UTSCALE
            vdplim=min(vdplim,.02)
            RSURFQ=RA(LCELL)+1./vdplim
    //       RSURFQ=RA(LCELL)+1./(UTSCALE*(Eff_dif+Eff_imp+Eff_int)*RBcor) // Total surface resistence
    //
//   limit this here to be consisten with the gocart routine, which bases this on Walcek et al. 1986
//
    //       RSURFQ=max(RSURFQ,50.)
            SUM0=SUM0+WGAUS(N)*(VSEDQ + 1./RSURFQ)  // Quadrature sum for 0 moment
            SUM3=SUM3+WGAUS(N)*(VSEDQ + 1./RSURFQ)*DQ**3  // Quadrature sum for 3rd moment
            ENDDO
            VDEP(LCELL, VDNNUC) = SUM0/sqrtpi  // normalize 0 moment vdep quadrature sum to sqrt(pi) (and number =1 per unit volume)
            VDEP(LCELL, VDMNUC) = SUM3/(sqrtpi*EXP((1.5*sqrt2*xxlsgn)**2)*DGNUC(LCELL)**3) //normalize 3 moment quad. sum to sqrt(pi) and 3rd moment analytic sum

// *** accumulation mode

            SUM0=0.
            SUM3=0.
            DO N=1,NGAUSdv
            DQ=DGACC(LCELL)*EXP(Y_GQ(N)*sqrt2*xxlsga)  // Diameter (m) at quadrature point
            KNQ=2.*XLM(LCELL)/DQ  // Knudsen number at quadrature point
            CUNQ=1.+KNQ*(1.257+.4*exp(-1.1/KNQ))  // Cunningham correction factor; Pruppacher and Klett (1980) Eq (12-16)
            VSEDQ=PDENSA(LCELL)*DCONST2*CUNQ*DQ*DQ  // Gravitational sedimentation velocity m/s
            SCQ=NU*DQ/DCONST1/CUNQ  // Schmidt number, Brownian diffusion parameter - Same as Binkowski and Shankar
            Eff_dif=SCQ**(-TWO3)    // Efficiency term for diffusion - Same as Binkowski and Shankar
            STQ=DCONST3*PDENSA(LCELL)*DQ**2  // Stokes number, Peters and Eiden (1992)
            Eff_imp=(STQ/(0.8+STQ))**2   // Efficiency term for impaction - Peters and Eiden (1992)
    //       Eff_int=0.3*DQ/(COLCTR_SMALD+DQ) // Slinn (1982) Interception term, 0.3 prefac insures .2 cm/s at .2 micron diam.
            Eff_int=(0.00116+0.0061*ZNTT(LCELL))*DQ/1.414E-7 // McKeen(2008) Intercptn term, val of .00421 @ ustr=0.475, diam=.1414 micrn, stable, needleleaf evergreen
            RBcor=1. // Rebound correction factor
            vdplim=UTSCALE*(Eff_dif+Eff_imp+Eff_int)*RBcor
            vdplim=min(vdplim,.02)
            RSURFQ=RA(LCELL)+1./vdplim
//       RSURFQ=RA(LCELL)+1./(UTSCALE*(Eff_dif+Eff_imp+Eff_int)*RBcor) // Total surface resistence
//
//   limit this here to be consisten with the gocart routine, which bases this on Walcek et al. 1986
//
//       RSURFQ=max(RSURFQ,50.)
        SUM0=SUM0+WGAUS(N)*(VSEDQ + 1./RSURFQ)  // Quadrature sum for 0 moment
        SUM3=SUM3+WGAUS(N)*(VSEDQ + 1./RSURFQ)*DQ**3  // Quadrature sum for 3rd moment
          IF(iprnt.eq.1)THEN
              print *,'N,Y_GQ,WGAUS,DQ,KNQ,CUNQ,VSEDQ,SCQ,STQ,RSURFQ'
              print *,N,Y_GQ(N),WGAUS(N),DQ,KNQ,CUNQ,VSEDQ,SCQ,STQ,RSURFQ
              print *,'N,Eff_dif,imp,int,SUM0,SUM3'
              print *,N,Eff_dif,Eff_imp,Eff_int,SUM0,SUM3
          endif
        ENDDO
        VDEP(LCELL, VDNACC) = SUM0/sqrtpi  // normalize 0 moment vdep quadrature sum to sqrt(pi) (and number =1 per unit volume)
        VDEP(LCELL, VDMACC) = SUM3/(sqrtpi*EXP((1.5*sqrt2*xxlsga)**2)*DGACC(LCELL)**3) //normalize 3 moment quad. sum to sqrt(pi) and 3rd moment analytic sum
        
// *** coarse mode 
        
        SUM0=0.
        SUM3=0.
        DO N=1,NGAUSdv
           DQ=DGCOR(LCELL)*EXP(Y_GQ(N)*sqrt2*xxlsgc)  // Diameter (m) at quadrature point
           KNQ=2.*XLM(LCELL)/DQ  // Knudsen number at quadrature point
           CUNQ=1.+KNQ*(1.257+.4*exp(-1.1/KNQ))  // Cunningham correction factor; Pruppacher and Klett (1980) Eq (12-16)
           VSEDQ=PDENSC(LCELL)*DCONST2*CUNQ*DQ*DQ  // Gravitational sedimentation velocity m/s
           SCQ=NU*DQ/DCONST1/CUNQ  // Schmidt number, Brownian diffusion parameter - Same as Binkowski and Shankar
           Eff_dif=SCQ**(-TWO3)    // Efficiency term for diffusion - Same as Binkowski and Shankar
           STQ=DCONST3*PDENSC(LCELL)*DQ**2  // Stokes number, Peters and Eiden (1992)
           Eff_imp=(STQ/(0.8+STQ))**2   // Efficiency term for impaction - Peters and Eiden (1992)
//          Eff_int=0.3*DQ/(COLCTR_SMALD+DQ) // Slinn (1982) Interception term, 0.3 prefac insures .2 cm/s at .2 micron diam.
           Eff_int=(0.00116+0.0061*ZNTT(LCELL))*DQ/1.414E-7 // McKeen(2008) Interception term, val of .00421 @ ustr=0.475, diam=.1414 micrn, stable, needleleaf evergreen
           EFF_int=min(1.,EFF_int)
           RBcor=exp(-2.0*(STQ**0.5)) // Rebound correction factor used in Slinn (1982)
           vdplim=UTSCALE*(Eff_dif+Eff_imp+Eff_int)*RBcor
           vdplim=min(vdplim,.02)
           RSURFQ=RA(LCELL)+1./vdplim
//       RSURFQ=RA(LCELL)+1./(UTSCALE*(Eff_dif+Eff_imp+Eff_int)*RBcor) // Total surface resistence
//
//   limit this here to be consisten with the gocart routine, which bases this on Walcek et al. 1986
//
//       RSURFQ=max(RSURFQ,50.)
           SUM0=SUM0+WGAUS(N)*(VSEDQ + 1./RSURFQ)  // Quadrature sum for 0 moment
           SUM3=SUM3+WGAUS(N)*(VSEDQ + 1./RSURFQ)*DQ**3  // Quadrature sum for 3rd moment
        ENDDO
            VDEP(LCELL, VDNCOR) = SUM0/sqrtpi  // normalize 0 moment vdep quadrature sum to sqrt(pi) (and number =1 per unit volume)
            VDEP(LCELL, VDMCOR) = SUM3/(sqrtpi*EXP((1.5*sqrt2*xxlsgc)**2)*DGCOR(LCELL)**3) //normalize 3 moment quad. sum to sqrt(pi) and 3rd moment analytic sum
        END DO
             
        ENDIF  // ENDOF LAYER = 1 test
        
// *** Calculate gravitational sedimentation velocities for all layers - as in Binkowski and Shankar (1995)

         DO LCELL = 1, NUMCELLS
         
            DCONST2 = GRAV / ( 18.0 * AMU(LCELL) )
            DCONST3N = DCONST2 * PDENSN(LCELL) * DGNUC( LCELL )**2
            DCONST3A = DCONST2 * PDENSA(LCELL) * DGACC( LCELL )**2
            DCONST3C = DCONST2 * PDENSC(LCELL) * DGCOR( LCELL )**2
               
// *** nucleation mode number and mass sedimentation velociticies
            VSED( LCELL, VSNNUC) = DCONST3N                         &
               * ( ESN16 + BHAT * KNNUC( LCELL ) * ESN04 )
            VSED( LCELL, VSMNUC) = DCONST3N                         &
               * (ESN64 + BHAT * KNNUC( LCELL ) * ESN28 )
        
// *** accumulation mode number and mass sedimentation velociticies
            VSED( LCELL, VSNACC) = DCONST3A                          &
              * ( ESA16 + BHAT * KNACC( LCELL ) * ESA04 )
            VSED( LCELL, VSMACC) = DCONST3A                          &
              * ( ESA64 + BHAT * KNACC( LCELL ) * ESA28 )

// *** coarse mode number and mass sedimentation velociticies
            VSED( LCELL, VSNCOR) = DCONST3C                          &
              * ( ESC16 + BHAT * KNCOR( LCELL ) * ESC04 )
            VSED( LCELL, VSMCOR) = DCONST3C                          &
              * ( ESC64 + BHAT * KNCOR( LCELL ) * ESC28 )
         END DO
END SUBROUTINE VDVG_2
